#   aws-cwa-metrics

AWS CloudWatchAgent Metrics Monitor

##  Installation

go get -u github.com/slatunje/aws-cwa-metrics

##  Deploy

On local machine    
    
```bash
with -ip vit-nonprod

aws s3 cp ${s3_source} ${s3_destination} --recursive
``` 

On ec2 instance - build the binary

```bash
wget https://s3.eu-west-1.amazonaws.com/${s3_destination}/1.0.0/linux/cwa-metrics

chmod +x cwa-metrics

mv ./cwa-metrics /usr/bin/cwa-metrics

cwa-metrics \
--mem \
--swap \
--disk \
--network \
--docker \
--region eu-west-1 \
--interval 1 \
--namespace CoreOS
``` 

On ec2 instance - create the service

```bash


```

##  Sources

-   https://github.com/shirou/gopsutil
-   http://www.blog.labouardy.com/publish-custom-metrics-aws-cloudwatch

