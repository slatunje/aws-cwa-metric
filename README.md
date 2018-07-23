#   aws-cwa-metrics

AWS CloudWatchAgent Metrics Monitor

##  Installation

go get -u github.com/slatunje/aws-cwa-metric

##  Deploy

On local machine    
    
```bash
with -ip vit-nonprod

aws s3 cp ${s3_source} ${s3_destination} --recursive
``` 

On ec2 instance - build the binary

```bash
wget https://s3.eu-west-1.amazonaws.com/${s3_destination}/1.0.0/linux/cwa-metric

chmod +x cwametrics

mv ./cwametrics /usr/bin/cwametrics

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

[SystemD Example](doc/unit.md)

##  Sources

-   https://github.com/shirou/gopsutil

##  Similar Repos

-   https://github.com/advantageous/metricsd
-   https://github.com/advantageous/systemd-cloud-watch
-   https://github.com/advantageous/journald-cloudwatch-logs  

-   https://github.com/saymedia/journald-cloudwatch-logs



