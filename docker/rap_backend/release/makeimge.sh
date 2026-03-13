#!/bin/bash 

IMAGETAG=$(date +%Y%m%d%H%M%S)_test
IMAGENAME=rap_backend
PROGRAMNAME=$0
REPOSITORY="airudder-registry.cn-hongkong.cr.aliyuncs.com/airudder"
BUILDIMAGE="airudder-registry.cn-hongkong.cr.aliyuncs.com/airudder/voip-go116-project-build:u18_plus"


WORKPATH=$(pwd)
echo $WORKPATH

function git_clone()
{
    if [ -d ./vans ]
    then
        rm -rf ./vans
    fi

    #git clone git@gitlab.corp.cootek.com:voip_infra/vans.git
    #cd ./vans && git submodule update --init --recursive && cd ..
}

function git_releasetag()
{
    cd $WORKPATH && \
    IMAGE="$REPOSITORY/$IMAGENAME:$IMAGETAG" && \
    git tag -a "v${IMAGETAG}" -m "version ${IMAGETAG}; docker pull ${IMAGE}" && \
    git push origin "v${IMAGETAG}"
}

function pull_build() 
{
    if [ -z "$ALIYUN_REGISTRY_USERNAME" ] || [ -z "$ALIYUN_REGISTRY_PASSWORD" ]; then
        echo "ALIYUN_REGISTRY_USERNAME or ALIYUN_REGISTRY_PASSWORD is not set"
        return 1
    fi
    printf '%s' "$ALIYUN_REGISTRY_PASSWORD" | docker login --username="$ALIYUN_REGISTRY_USERNAME" --password-stdin airudder-registry.cn-hongkong.cr.aliyuncs.com
    docker pull $BUILDIMAGE
}

function build_code2()
{
    CODEPATH=$WORKPATH/../../../../rap_backend
        echo $CODEPATH
    docker run -it --rm=true -v $CODEPATH:/go-project \
        $BUILDIMAGE /bin/bash

}

function build_code()
{
    CODEPATH=$WORKPATH/../../../../rap_backend
        echo $CODEPATH
    docker run -it --rm=true -v $CODEPATH:/go-project \
        $BUILDIMAGE /bin/bash /go-project/build.sh rap_backend && \
    cp $CODEPATH/bin/rap_backend $WORKPATH/rap_backend

}

function docker_build()
{
    docker build -t $REPOSITORY/$IMAGENAME:$IMAGETAG -f Dockerfile.release .    && \
    docker tag $REPOSITORY/$IMAGENAME:$IMAGETAG  $REPOSITORY/$IMAGENAME:latest  && \
    echo "build finish $REPOSITORY/$IMAGENAME:$IMAGETAG"                        && \
    rm -rf $WORKPATH/rap_backend
}

function docker_push()
{
    docker push $REPOSITORY/$IMAGENAME:$IMAGETAG   && \
    docker push $REPOSITORY/$IMAGENAME:latest      && \
    echo "push finish $REPOSITORY/$IMAGENAME:$IMAGETAG"
}

#git_clone && pull_build && \
#     build_code && docker_build && docker_push && git_releasetag
#build_code2

pull_build && build_code && docker_build && docker_push
