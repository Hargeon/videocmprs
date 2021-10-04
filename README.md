# videocmprs - video compression service

## Description
The service expose a RESTful API to
downsample and compress the video with the ratio and bitrate specified by the user. User
able to see requests history and status (e.g. queued, processing, done) and
download the original video and the processed one.

## Create migration
```bash
goose -dir db/migrations/common create users sql
```

## Migrate
```bash
goose -dir db/migrations/common postgres "user=postgres dbname=... sslmode=disable host=... port=... password=..." up
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
videocmprs is released under the [MIT License](https://github.com/Hargeon/videocmprs/blob/master/LICENSE).