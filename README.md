# Password manager GophKeeper

GophKeeper is a client-server system that allows the user to safely and securely store logins, passwords, binary data and other private information.

## General requirements

#### Server side

- user registration, authentication and authorization;

- storing the private data;

- data synchronization between several authorized clients of the same owner;

- transfer the requested private data to the owner.

#### Client side

- authentication and authorization for users on the remote server;

- access to private data by request.

#### Functions implementation of which is up to the implementer's discretion

- server-side or client-side data creating, editing, and deleting;

- format of new user registration;

- selection of storage and data storage format;

- ensuring security of data transmission and storage;

- client-server interaction protocol;

- mechanisms for user authentication and authorization of access to information.

#### Additional requirements

- Client chould be distributed as a CLI application with the ability to run on Windows, Linux, and Mac OS platforms;

- Client should allow the user to get information about the version and build date of the client binary file.

## Quick start

To display all possible commands:

`make help`

To run the server:

`make server/run`

To run the client:

`make client/run`