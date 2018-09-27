package auth

import (
	"errors"
	"spyc_backend/app"
	"spyc_backend/app/core"
	"spyc_backend/app/models"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/revel/revel"
)

// Const for auth
const (
	SigningKey = "$3cr3t_k3y!"
)

// Claims Estructura para el manejo de la seccion Payload en el token.
type Claims struct {
	jwt.StandardClaims
	ID     string `json:"id"`
	Action string `json:"action"`
}

// GameClaims is the info used in a game token
type GameClaims struct {
	jwt.StandardClaims
	UserID             string `json:"ui"`
	StepSubscriptionID string `json:"ssi"`
	Action             string `json:"action"`
}

type ChallengeClaims struct {
	jwt.StandardClaims
	UserID      string `json:"ui"`
	ChallengeID string `json:"ci"`
	Action      string `json:"action"`
}

// GenerateToken Funcion que genera un token de authenticacion para un usuario.
func GenerateToken(id string, action string) (token string, err error) {
	// Tiempo de expiracion del Token, 1 semana
	// expireToken := time.Now().Add(time.Hour * 24 * 7).Unix()

	// Se crea la estructura del Payload.
	claims := Claims{
		ID:     id,
		Action: action,
	}

	// Se genera el token y se firma.
	tokenKey := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenKey.SignedString([]byte(SigningKey))

	return
}

// GenerateGameToken creates the token required for a game
func GenerateGameToken(idSubscription, idUser, action string) (token string, err error) {
	gClaims := GameClaims{
		UserID:             idUser,
		StepSubscriptionID: idSubscription,
		Action:             action,
	}

	tokenKey := jwt.NewWithClaims(jwt.SigningMethodHS256, gClaims)
	token, err = tokenKey.SignedString([]byte(SigningKey))
	return
}

// GenerateChallengeToken generates token for challenges
func GenerateChallengeToken(idChallenge, idUser, action string) (token string, err error) {
	gClaims := ChallengeClaims{
		UserID:      idUser,
		ChallengeID: idChallenge,
		Action:      action,
	}

	tokenKey := jwt.NewWithClaims(jwt.SigningMethodHS256, gClaims)
	token, err = tokenKey.SignedString([]byte(SigningKey))
	return
}

// VerifyToken ...
func VerifyToken(token, action string) (*Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected siging method")
		}

		return []byte(SigningKey), nil
	})
	if err != nil {
		return &Claims{}, err
	}

	// Verificamos la firma del Token, si el valida obtenemos el id del usurio.
	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		if claims.Action != action {
			return claims, errors.New("Action not match")
		}

		return claims, err
	}

	return &Claims{}, errors.New("Invalid token")
}

// VerifyGameToken validates the one-use token
func VerifyGameToken(gameToken, action string) (err error) {
	parsedToken, err := jwt.ParseWithClaims(gameToken, &GameClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected siging method")
		}
		return []byte(SigningKey), nil
	})
	if err != nil {
		return
	}

	// Verificamos la firma del Token, si es valida obtenemos el id del usurio y el id del step
	if GameClaims, ok := parsedToken.Claims.(*GameClaims); ok && parsedToken.Valid {
		if GameClaims.Action != action {
			return errors.New("Action not match")
		}
		return
	}
	return errors.New("Invalid token")
}

// AuthenticateByToken funcion para obtener un usuaio con el token.
func AuthenticateByToken(req *revel.Request) (user models.User, err error) {
	authHeader := req.Header.Get("Authorization")

	if !strings.Contains(authHeader, core.Bearer) {
		err = errors.New("Missing token")
		return
	}

	token := strings.Split(authHeader, core.Bearer)[1]

	// Se intenta verificar el Token (algoritmo, y payload)
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected siging method")
		}

		return []byte(SigningKey), nil
	})
	if err != nil {
		// Retornamos el error al no porder comprovar el Token.
		return user, err
	}

	// Verificamos la firma del Token, si el valida obtenemos el id del usurio.
	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		// Buscamos a el usuario en la base de datos por ID
		if User, ok := app.Mapper.GetModel(&models.User{}); ok {
			if err = User.Find(claims.ID).Exec(&user); err != nil {
				return user, err
			}
		}
	}

	return user, nil
}

// AuthenticateBySession funcion para obtener un usuaio con el token.
func AuthenticateBySession(token string) (user models.User, err error) {
	// Se intenta verificar el Token (algoritmo, y payload)
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected siging method")
		}

		return []byte(SigningKey), nil
	})
	if err != nil {
		return user, err
	}

	// Verificamos la firma del Token, si el valida obtenemos el id del usurio.
	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		// Buscamos a el usuario en la base de datos por ID
		User, _ := app.Mapper.GetModel(&models.User{})
		err = User.Find(claims.ID).Exec(&user)
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

// AuthenticateGameToken Validates token with the given ation
// returns subscription(TargetUser) and token user
func AuthenticateGameToken(req *revel.Request, action string) (subscriptionstep models.TargetUser, user models.User, err error) {

	authHeader := req.Header.Get("Authorization")

	if !strings.Contains(authHeader, core.Bearer) {
		err = errors.New("Missing token")
		return
	}

	token := strings.Split(authHeader, core.Bearer)[1]

	// Se intenta verificar el Token (algoritmo, y payload)
	parsedToken, err := jwt.ParseWithClaims(token, &GameClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected siging method")
		}
		return []byte(SigningKey), nil
	})
	if err != nil {
		return subscriptionstep, user, err
	}

	// Verificamos la firma del Token, si el valida obtenemos el id del usurio.
	if claims, ok := parsedToken.Claims.(*GameClaims); ok && parsedToken.Valid {

		// Validamos la action a comparar con la del token
		if claims.Action != action {
			return subscriptionstep, user, errors.New("Invalid validation")
		}

		// Find subscription based on stepSubscriptionID from token
		TargetUser, _ := app.Mapper.GetModel(&subscriptionstep)
		err = TargetUser.Find(claims.StepSubscriptionID).Exec(&subscriptionstep)
		if err != nil {
			return subscriptionstep, user, err
		}

		// Buscamos a el usuario en la base de datos por ID
		User, _ := app.Mapper.GetModel(&models.User{})

		err = User.Find(claims.UserID).Exec(&user)
		if err != nil {
			return subscriptionstep, user, err
		}
	}
	return subscriptionstep, user, err
}

// AuthenticateChallengeToken ...
func AuthenticateChallengeToken(req *revel.Request, action string, status string) (challenge models.Challenge, user models.User, err error) {
	authHeader := req.Header.Get("Authorization")

	if !strings.Contains(authHeader, core.Bearer) {
		err = errors.New("Missing token")
		return
	}

	token := strings.Split(authHeader, core.Bearer)[1]

	// Se intenta verificar el Token (algoritmo, y payload)
	parsedToken, err := jwt.ParseWithClaims(token, &ChallengeClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected siging method")
		}

		return []byte(SigningKey), nil
	})
	if err != nil {
		return challenge, user, err
	}

	// Verificamos la firma del Token, si el valida obtenemos el id del usurio.
	if claims, ok := parsedToken.Claims.(*ChallengeClaims); ok && parsedToken.Valid {

		// Validamos la action a comparar con la del token
		if claims.Action != action {
			return challenge, user, errors.New("Invalid validation")
		}

		// Buscamos a el usuario en la base de datos por ID
		User, _ := app.Mapper.GetModel(&models.User{})
		if err = User.Find(claims.UserID).Exec(&user); err != nil {
			return challenge, user, err
		}

		var params = []bson.M{
			bson.M{"_id": bson.ObjectIdHex(claims.ChallengeID)},
			bson.M{"status.name": status},
		}
		Challenge, _ := app.Mapper.GetModel(&challenge)
		if err = Challenge.FindWithOperator("$and", params).Exec(&challenge); err != nil {
			return challenge, user, err
		}

	}
	return challenge, user, nil
}

// GetHeaderCoordinates verifies if user location in
// Answer validation exists and has 2 arguments
// latitud <= +-90 && longitud <= +-180
func GetHeaderCoordinates(req *revel.Request) (geo models.Geo, err error) {
	authHeader := req.Header.Get("Geo")

	if !strings.Contains(authHeader, core.UserLocation) {
		err = errors.New("Missing location")
		return geo, err
	}

	coordsString := strings.Split(authHeader, core.UserLocation)[1]

	coords := strings.Split(strings.Trim(coordsString, " "), ",")

	if len(coords) <= 1 {
		revel.ERROR.Print("ERROR COORDS Invalid Coordinates " + authHeader)
		return geo, errors.New("Invalid header")
	}

	lng, errLng := strconv.ParseFloat(coords[0], 64)
	lat, errLat := strconv.ParseFloat(coords[1], 64)

	if errLat != nil || errLng != nil {
		return geo, errLat
	}

	if lng > 180 || lng < -180 || lat > 90 || lat < -90 {
		return geo, errors.New("Invalid Coordinates")
	}

	geo.Coordinates = []float64{lng, lat}

	return geo, nil

}

// GetHeaderLanguage gets the language from the header "Accept-Language""
func GetHeaderLanguage(req *revel.Request) (lang string) {
	lang = req.Header.Get("Accept-Language")
	return lang
}
