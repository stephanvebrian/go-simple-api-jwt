## go-simple-api-jwt

act as an proxy for http://13.212.226.116:8000/docs/

p.s: please check the /register endpoint from above api, looks like its broken when trying to upload the file image, you can check using below curl (adjust the image path),
if we trace the error its actually caused by this one https://stackoverflow.com/questions/67089855/str-object-has-no-attribute-decode-on-djangorestframework-simplejwt

```curl
curl --location 'http://13.212.226.116:8000/api/register' \
--header 'Content-Type: multipart/form-data' \
--form 'username="admin12345@gmail.com"' \
--form 'password="123"' \
--form 'first_name="First Admin"' \
--form 'last_name="Last Admin"' \
--form 'telephone="+62888888888888"' \
--form 'profile_image=@"/Users/stephan.vebrian/Downloads/test-1.jpg"' \
--form 'address="Jln Palmerah 53"' \
--form 'city="Jakarta Barat"' \
--form 'province="DKI Jakarta"' \
--form 'country="Indonesia"' \
--form 'email="admin12345@gmail.com"'
```
