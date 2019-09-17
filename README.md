# Go Tensorflow (GOTF)

Serving the Tensorflow Resnet Model using Docker and:
- Basic Python
- Compiled Go
- [Tensorflow Serving](https://github.com/tensorflow/serving)

Performance tested using [Vegeta](https://github.com/tsenart/vegeta)

### Serving

Change into the example you'd like to run. 

This is a resnet model.

Build and run docker images for Basic and Compiled Go servers

```
$ docker build -f docker/base.Dockerfile . -t base ; docker run -p 8501:8080 -it base:latest

$ docker build -f docker/compiled.Dockerfile . -t compiled ; docker run -p 8501:8080 -it compiled:latest
```

Tensorflow serving 
```
docker run -p 8501:8501 \
  --mount type=bind,source=/Users/alicecheung/Downloads/gotf/model,target=/models/resnet_model \
  -e MODEL_NAME=resnet_model -t tensorflow/serving
```

### Performance Testing

Run test in testing directory

```
$ go run test.go
```


