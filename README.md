<h1 align="center"> CSV -upload  </h1>
<h2> POSTMAN COLLECTION TO TEST THE API </h2>

-> https://www.getpostman.com/collections/eb2a228cd6ea8c04021d

<h2> Video demo of the API </h2>

-> https://www.youtube.com/watch?v=TT7_cCSFHU0


<h2> Route endpoints </h2>

<h4> Local URL -> localhost:8080 </h4>

- `{{url}}/users`        :  Post request to add user to the database.

- `{{url}}/users/:id`      :  Get request to get the user profile based on ID.

- `{{url}}/posts`     :  Post request to create a POST under an AuthorID. 

- `{{url}}/posts/p1`      :  GET request to get the Post based on its ID.

- `{{url}}/posts/users/:id?page={{page number}}`     :  Get request to get the Posts of a user with pagination(3 posts in each page).






<h3>How to run the server locally</h3>

- Run `https://github.com/C-Harshul/Insta-API.git` on your terminal

- Go to the terminal and run `go run main.go`

The server should be online on port 8080


