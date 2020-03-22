# s3reverse

<img src="https://img.shields.io/github/languages/top/hahwul/s3reverse"> <img src="https://img.shields.io/github/license/hahwul/s3reverse.svg"> <a href="https://twitter.com/intent/follow?screen_name=hahwul"><img src="https://img.shields.io/twitter/follow/hahwul?style=flat-square"></a>

## Install
```cassandraql
$ go get -u github.com/hahwul/s3reverse
```
## Usage
### Input options
Basic Usage
```cassandraql

8""""8 eeee       8"""8  8"""" 88   8 8"""" 8"""8  8""""8 8""""
8         8       8   8  8     88   8 8     8   8  8      8
8eeeee    8       8eee8e 8eeee 88  e8 8eeee 8eee8e 8eeeee 8eeee
    88 eee8  eeee 88   8 88    "8  8  88    88   8     88 88
e   88    88      88   8 88     8  8  88    88   8 e   88 88
8eee88 eee88      88   8 88eee  8ee8  88eee 88   8 8eee88 88eee

by @hahwul

Usage of ./s3reverse:
  -iL string
    	input List
  -oA string
    	Write output in Array format (optional)
  -oN string
    	Write output in Normal format (optional)
  -tN
    	to name
  -tP
    	to path-style
  -tS
    	to s3 url
  -tV
    	to virtual-hosted-style
  -verify
    	testing bucket(acl,takeover)
```
Using from file
```cassandraql
$ s3reverse -iL sample -tN
udemy-web-upload-transitional
github-cloud
github-production-repository-file-5c1aeb
github-production-upload-manifest-file-7fdce7
github-production-user-asset-6210df
github-education-web
github-jobs
s3-us-west-2.amazonaws.com
optimizely
app-usa-modeast-prod-a01239f
doc
swipely-merchant-assets
adslfjasldfkjasldkfjalsdfkajsljasldf
cbphotovideo
cbphotovideo-eu
public.chaturbate.com
wowdvr
cbvideoupload
testbuckettesttest
```
Using from pipeline
```cassandraql
$ cat sample | s3reverse -tN
udemy-web-upload-transitional
github-cloud
github-production-repository-file-5c1aeb
github-production-upload-manifest-file-7fdce7
github-production-user-asset-6210df
github-education-web
github-jobs
s3-us-west-2.amazonaws.com
optimizely
app-usa-modeast-prod-a01239f
doc
swipely-merchant-assets
adslfjasldfkjasldkfjalsdfkajsljasldf
cbphotovideo
cbphotovideo-eu
public.chaturbate.com
wowdvr
cbvideoupload
testbuckettesttest
```

### Output options
to Name
```cassandraql
$ s3reverse -iL sample -tN
udemy-web-upload-transitional
github-cloud
github-production-repository-file-5c1aeb
github-production-upload-manifest-file-7fdce7
... snip ...
```
to Path Style
```cassandraql
$ s3reverse -iL sample -tP
https://s3.amazonaws.com/udemy-web-upload-transitional
https://s3.amazonaws.com/github-cloud
https://s3.amazonaws.com/github-production-repository-file-5c1aeb
... snip ...
```
to Virtual Hosted Style
```cassandraql
$ s3reverse -iL sample -tV
udemy-web-upload-transitional.s3.amazonaws.com
github-cloud.s3.amazonaws.com
github-production-repository-file-5c1aeb.s3.amazonaws.com
github-production-upload-manifest-file-7fdce7.s3.amazonaws.com
github-production-user-asset-6210df.s3.amazonaws.com
... snip ...
```

### Verify mode
```cassandraql
$ s3reverse -iL sample -verify
[NoSuchBucket] adslfjasldfkjasldkfjalsdfkajsljasldf
[PublicAccessDenied] github-production-user-asset-6210df
[PublicAccessDenied] github-jobs
[PublicAccessDenied] public.chaturbate.com
[PublicAccessDenied] github-education-web
[PublicAccessDenied] github-production-repository-file-5c1aeb
[PublicAccessDenied] testbuckettesttest
[PublicAccessDenied] app-usa-modeast-prod-a01239f
[PublicAccessGranted] cbphotovideo-eu
[PublicAccessDenied] swipely-merchant-assets
[PublicAccessDenied] optimizely
[PublicAccessDenied] wowdvr
[PublicAccessGranted] s3-us-west-2.amazonaws.com
[PublicAccessDenied] cbphotovideo
[PublicAccessDenied] cbvideoupload
[PublicAccessDenied] github-production-upload-manifest-file-7fdce7
[PublicAccessDenied] doc
[PublicAccessDenied] udemy-web-upload-transitional
[PublicAccessDenied] github-cloud
```

## Case study
Pipelining `meg`, `s3reverse`, `gf` , `s3scanner` for Find S3 Misconfiguration.
```cassandraql
$ meg -d 1000 -v / ; cd out ; gf s3-buckets | s3reverse -tN > buckets ; s3scanner buckets
```

Find S3 bucket takeover
```cassandraql
$ meg -d 1000 -v / ; cd out ; gf s3-buckets | s3reverse -verify | grep NoSuchBucket > takeovers
```
