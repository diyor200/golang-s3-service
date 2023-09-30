<h1>S3 service for uploading files to AWS S3 with authentication system</h1>
1. First rename <i>.env.example`</i> file to <i>.env</i> and enter your data to it corresponding order.
<br>
2. Run <code>docker compose up</code>
<br>
<h3>Endpoints:</h3>
<h4>1. <code>localhost:8080/auth/sign-up</code> send request it with your json data like this</h4>
<br>
<code>{ "username": "example_user1", "password":"examplePasword1@"}</code> 
After successfully auth returns <code>id</code> of the user.
<br>
<image src="images/sign-up.png"></image>
<br>
<h4>2. <code>localhost:8080/auth/sign-in</code> send your json data:</h4><br>
<code>{ "username": "example_user1", "password":"examplePasword1@"}</code>. After successfully
sign-in it returs <cod>jwt</cod> token like this: <br>
<image src="images/sign-in.png"></image><br>
<h4>3. <code>localhost:8080/api/v1/file</code> send your image file with <code>file</code> key in header. System accepts only image files.
To send your image you need to enter jwt token. After that You can send image. Image size must be no more than 8MBs.
After successfully sent file system return link of this file in the cloud:</h4>
<img src="images/upload1.png">
<img src="images/upload2.png">.
<br>
<a href="https://www.buymeacoffee.com/diyorbekabdulaxatov" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>