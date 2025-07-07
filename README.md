# Services Golang Architecture

This project follows key architectural principles to maintain clean, maintainable, and scalable code.

## Core Design Philosophy

### 1. Package Boundaries
- Each package acts as a firewall/boundary in the system
- Packages organize APIs rather than just code
- APIs are layered vertically within their domain

### 2. Package Purpose
- Packages provide specific business purposes
- Handle necessary data transformations
- Focus on "what they do" rather than "what they contain"

### 3. Type System
- Each package defines its own type system
- Types represent both incoming and outgoing data
- Prefer concrete types (users, products, sales) for clarity

### 4. Polymorphism Usage
- Use interfaces for behavior-based operations
- Runtime polymorphism: achieved through interfaces
- Static polymorphism: achieved through generics

### 5. Interface Guidelines
- Use interfaces as input types when needed
- Avoid interfaces as return types
- Let callers handle data decoupling
- Return concrete type pointers for clarity 

### 6. Application Layer
- App layer must be agnostic to the protocol layer
- Business logic should be independent of delivery mechanisms
- Enables flexibility to change protocols without affecting core functionality
- Promotes clean separation of concerns and testability 


## Authorization

 
We will associate a Private Key with our system and create a signature associated with our data. Every time our data changes, the signature will be different due to the relationship between data and the private key. The data and the signature will be converted into a Base64 Encoded Token (not encrypted, but encoded).
We will generate a data payload, send it to the server, sign it with our private key, produce a token with information encoded in base64, and pass it to the user. We will instruct the user to pass this token as an authorization header. When the token arrives, we will verify if it is the token we produced and if it was signed with our private key. If not, we will take appropriate action.

Upon completion of these steps, we will proceed with the authorization process. We will retrieve claims from the token to ascertain the user’s identity and delegate the authorization process to the Open Policy Agent system (https://www.openpolicyagent.org/). Instead of implementing the logic in GoLang, we will utilize a separate server called Rigo to execute the functionality. The general approach for this step involves defining middleware that establishes the necessary claims for the user to perform the endpoint operation.

In a production system, I typically advocate against developing our own JWT system. Instead, I recommend using a reputable third-party service that manages user authentication, token generation, and user management. Our system should simply load the public keys for the private keys (pairs) and utilize those public keys for authentication purposes. This approach avoids the need to handle private keys directly, which poses a security concern that we prefer to avoid. However, for educational purposes, we will implement our own system.
The JWT protocol serves as a fundamental framework for our system, but we will need to add additional logic to enhance its functionality.

However, the JWT protocol authentication process alone is insufficient. Consider the scenario where we wish to restrict access to a user who possesses a valid token. While we could argue that the token has a Time-to-Live (TTL), we must wait until that expiration date. As an alternative and drastic measure, we could invalidate the private key and, consequently, the token. However, this approach would result in the non-effective invalidation of every token, which is not appropriate. Therefore, we require an additional step to address this issue. The JWT protocol authentication process is commonly referred to as Level One authentication, but it falls short of providing comprehensive security. After successful Level One authentication, we, after extracting the user ID from the token, typically proceed to the database verification. Subsequently, we verify whether the user is still enabled within the system.	