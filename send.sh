#!/usr/bin/env sh

function send() {
   depId=MSCLNckJ9ZHzNZYmBT69YbbgPwMik1yz7P3S
   depId=MSCL00000000000000000000000000000001
   productType=MDCORE
   productUrl=localhost:9093
   apikey=aaaaaaaaaaaaaaaaaaaaaaaa
   curl -sSL -w 'status: %{http_code}\n' -X POST $productUrl'/mdcore/integration/console/'$depId'/report/scan' \
        -H 'content-type: multipart/form-data' \
        -H 'api-key: '$apikey \
        -F "file01=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/01a92ff0f3e64e13a14ca08c824902d4.json" \
        -F "file02=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/01fdbc2c4d574d8c8a86e40be6d56c91.json" \
        -F "file03=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0a0586ad64a041209c6f1f29a4d9be29.json" \
        -F "file04=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0a9ddb1a1263440bb748756944a2c70d.json" \
        -F "file05=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0aa2c25b0f0f4c549d8a986dfecb66e0.json" \
        -F "file06=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0aa63dafff8c4f7da3bbfd915cf3010a.json" \
        -F "file07=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0ad7595041a4469b84d9c4bb20b4c6c2.json" \
        -F "file08=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0b04c61336454b69beb6f2b8bc60c8c3.json" \
        -F "file09=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0bae971e023a465dbae1e5f3192d403f.json" \
        -F "file10=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0c1426a7df884ecdb84c2be9478f49b8.json" \
        -F "file11=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0c3f2196144945a988a6cd90188ecff7.json" \
        -F "file12=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0c859a9c06f04efa8bcc892ae8922698.json" \
        -F "file13=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0cadac973db74c76b44e14ab1d61fdf3.json" \
        -F "file14=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0cce207749834b7aa4c842386cabedd8.json" \
        -F "file15=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0d4e4c4c11df4631bac80a19c028ae0e.json" \
        -F "file16=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0daa34ac01b0408eb162b996e8cd379d.json" \
        -F "file17=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0e9a162d541c4045835227ae7d8f7ed0.json" \
        -F "file18=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0ec12fae8fb448a79ce10a7521df1062.json" \
        -F "file19=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0ed413136cd14af7b3aa911b76239848.json" \
        -F "file20=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0ed47502c6af43c88a709aabf69b1a1e.json" \
        -F "file21=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0eee3abdf6bf43cf81d64cf43ab0b14c.json" \
        -F "file22=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0eff55da5b9d44168f4774885b045818.json" \
        -F "file23=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0f11b39ee42c45fdbe2605847a223ec5.json" \
        -F "file24=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0f3e12ff1d4f448983e751c5e669f62a.json" \
        -F "file25=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0f46e793001d48b383de856438ff5b01.json" \
        -F "file26=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0f8a0d30e34a4566b0307abab523df7b.json" \
        -F "file27=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0fcb2854ea1e4b4ab6e7a04e74fa3334.json" \
        -F "file28=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/0fe331b8e14a430abfde88d4716ab336.json" \
        -F "file29=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/1a5898f8b74f480ab518758605e437cd.json" \
        -F "file30=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0/1a9e1664ce9e4a6a8b6f9bb2aa9d8478.json" \
        -F "file31=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629/e777042cd0904c5cbd07515903066cf0.json" \
        -F "file32=@/home/minh/Documents/dev/core_req_nested/f336c57e7c4045db9d8771d29872f629.json"
}

for i in $(seq 1 1); do
   send
done
