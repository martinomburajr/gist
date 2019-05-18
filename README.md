# Gist

Gist is an application that takes in a file of any kind and uploads it to gist.github.com under your username. 
Metadata for the file is supplied within the file itself. API information should be supplied on initialization. 
This application written in Golang.

## How it works
The application binary is ran using the term: 
    
    gist

### 1. App Initialization
The first time you start the application. It will open a browser and request you to log into GitHub. This process 
follows the typical OAuth2.0 flow and grants the client (this application) access to specified scopes. These **scopes** 
are the following:

    1. gist
    
### 2. "Gisting"
Once signed in, you can "gist" a file or set of files of your choice. 

#### 2.1 "Gisting" a File
##### 2.1.1 Gist Flags
Given a file, `myfile.txt`, to declare it `gistable` i.e able to be sent to gist.github.com to your account, you need
 to supply metadata to the `gist` client.
 
###### 2.1.1.1 Gist an entire file
 To gist an entire file run the following command (once logged in). Assume the file is called `file.txt` and exists 
 within the `/home/me/code/` directory. Simply run the command:
 
 **absolute path:** `gist /home/me/code/file.txt -a "Martin Ombura Jr" -d "This is a file that has some text"`
 
 **relative path:** `gist file.txt -a "Martin Ombura Jr" -d "This is a file that has some text"`
 
## All Flags
    push : mirrors git's push command to upload the selected files and content to the server. The file or files must 
    be the last flags in the command
    -a : Specify a author of the gist
    -d : Description for the gist
    -n : Name of the gist when uploaded. This should simply be a file name e.g. main.go, upload.py etc.
 
## Contribute
Feel free to create issues/pull requests or fork the repo for your own usage!
