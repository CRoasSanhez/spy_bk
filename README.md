# Spy chatter backend

·## MongoDB Index

email, username unique index
+ db.users.createIndex({"user_name":1}, { unique: true, partialFilterExpression: { "user_name": { "$gt": ""}}})

### Dependencies

1. Change to '**/utilities**' directory.
```
cd utilities/
```

2. Add execution permission to '**dependencies.sh**' script.
```
chmod +x dependencies.sh
```
3. Run it.
```
./dependencies.sh
```

Done!


Cross compile permission denied error
ls $(go env GOROOT)/pkg && sudo chown -R $USER $(go env GOROOT)/pkg
