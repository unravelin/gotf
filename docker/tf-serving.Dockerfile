FROM tensorflow/serving:1.13.0

ENV MODEL_NAME=resnet_model

COPY model/ /models/resnet_model