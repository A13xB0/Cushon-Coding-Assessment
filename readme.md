# Coushin Interview Coding Assessment
## The task
### Description
Cushon already offers ISAs and Pensions to Employees of Companies (Employers) who have an existing arrangement with Cushon. Cushon would like to be able to offer ISA investments to retail (direct) customers who are not associated with an employer. Cushon would like to keep the functionality for retail ISA customers separate from it’s Employer based offering where practical.  
  
When customers invest into a Cushon ISA they should be able to select a single fund from a list of available options. Currently they will be restricted to selecting a single fund however in the future we would anticipate allowing selection of multiple options.  
  
Once the customer’s selection has been made, they should also be able to provide details of the amount they would like to invest. Given the customer has both made their selection and provided the amount the system should record these values and allow these details to be queried at a later date.

As a specific use case please consider a customer who wishes to deposit £25,000 into a Cushon ISA all into the Cushon Equities Fund

### Assessment
Please provide your solution to the above scenario in whatever form you feel is appropriate, using your preferred tools.  
Please spend the amount of time you feel appropriate to showcase your abilities and knowledge.  
Please be prepared to discuss during an interview:  
* What you have done and why.
* The specific decisions you made about your solution.
* Any assumptions you have made in the solution you have presented.
* Any enhancements you considered but decided not to cover.  

## My background
I am a Go Developer who leads, developes, and designs services within a client on premise environment and at the edge of their network. I do not do front-end services and a lot of my services is to do with data processing and pushing to ledgers and filestores. I decided as my skillset is backend I would set the scope of the assessment to be backend REST Endpoints for a front end to interface with. 
  
In this assessment I intended to show my knowledge in the following concepts:
* Common Go Patterns (Object Orientated Programming, Interfaces, composition, channels, signal handling, elastic search compatible logging, using packages).
* K8s (Not completed)
* Dockerfiles
* Authentication (JWT)
* Using Web Frameworks
* Microservices
* Managing a mono repo

## The result
### Disclaimer
Unfortunately I had a very short time frame with my current duties in my current job to complete these tasks having spent 3 days in London (2 overnight) last week and having a busy weekend, followed by working away on Monday I only had a very short time on the weekend to complete what I have. This means the following is missing:
* Untested code (but it compiles) - This code is untested due to a firewall rule expiring on my corporate laptop (my personal laptop is out of action) prevent me pulling from docker. This shouldn't have been an issue but it cropped up and the firewall team have yet to reply to my ticket. 
* No tests - I usually would do TDD but due to not having access to postgres in a nice docker environment I decided against this and just get on with the task at hand and write the tests after. Unfortunately the interview arrived a day before I expected due to availabilities.
* K8s - I was going to use kubernetes in this architecture for the microservices, unfortunately I never got around to it
* Pipeline - I usually use Go Sec, Go Vet, Golangci-lint, and sonarqube to detect any security vulnerabilities and linting issues in my code within my taskfile (basically like a makefile but more modern and simpler to use). But the interview crept up on me so I have not managed to do this. 
### Architecture
#### Technologies
I decided to use the following technologies and frameworks:
* Go 1.20.5 - This is what I usually develop microservices and native AWS applications in and is what I am proficient in.
* Gin - This is a personal preference in web framework for developing web APIs
* JWT - This allows for the use of bearer tokens for stateless authentication to other services once authenticated
* Taskfile - This is a more modern type of makefile which more resembles a buildspec file for building.
* Docker - For creating the containers to place into k8s
* K8s - For creating pods and adding scaling to the design
* Zap logger - Ubers logging in Go, a very capable and well liked in the Gopher community.

#### Microservices
I decided to go with a microservice architecture consisting of 3 services:
* Authentication - This is a REST API used for authentication by claiming JWT tokens and encapsulating the user information in the payload and signing the data to allow for stateless authentication between services. The downfall of JWT is revoking tokens can be difficult. A username and password is checked and then the JWT token is produced. The package which manages the JWT token is in a utility so the other services can act as validators.
* Customer - This is a REST API used for editing and modifying any customer information. This is intended to be used by a front end and provides only a small amount of example REST endpoints. 
* Investments - This is a REST API used for submitting investments to the SQL database. I personally have no idea how investment APIs work however in this example I have given I would presume some processor could process the investment rows.

The intention of these endpoints was to seperate out the services provided into smaller microservices which do very specific tasks. In my example they all use the same database which is relational and uses postgres, however in this design it is possible to switch to a NoSQL type database by forming relationships in the services themselves or denormalising by duplicating data between stores. This would allow great horizontal scaling with larger databases but is less strict with data which is a drawback.

The thought process was that this could split REST APIs into very specific tasks for example the following:
* https://authentication.api.coushin.com/api/v1/
* https://customer.api.coushin.com/api/v1/
* https://investment.api.coushin.com/api/v1/

This allows the seperation of the tasks and also if one REST Endpoint unfortunately fails the others will continue as normal.

### My improvements
The following I would either change or not do in the real world:
* Instead of flags I would create a config builder and use viper to allow the use of environment variables and config files for within containers. This simplifies running the services.
* Not use an interface across a database. This was only done to demonstrate the use of an interface as this is an important concept in Go. Doing this on a database creates added complexity which is not needed. This would be better used in something like a file handler to switch between S3, Samba, etc. 
* I would implement a database reconnect into the microservices, this does not exist in the database packages of Go and if there was a database connection interruption a nil pointer will probably occur.
* I would not write my own authetication service, there is a lot of risk in writing your own as this is the first thing malicious actors like to attack. I would instead rely in a commercial off the shelf service, self hosted or SaaS.
* If I had to write an authentication service I would not just hash the password using Blake2B, I would salt the password which prevents attacks using prehashed rainbow tables.
* I would add auditing using middleware into Gin.
* I would design the API using protobuf to allow for the use of GRPC for interservice communication if required aswell as using the protojson for REST APIs. Alternatively I could use swagger. 
* I would have written this using test driven development.
