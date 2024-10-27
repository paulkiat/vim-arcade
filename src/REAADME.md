.value;
    const password = req.body.password;
    const name = req.body.name;
    const email = req.body.email;
    const passwordHash = 
    
# Arena Server Protocol Documentation

## Introduction

 This document describes a protocol for communicating with an arena server. It is not complete, and there will be more content to this document as time goes on. If you have any questions or suggestions about the protocol or the documentation, please email <Paul Kiat> <<pk2brain@gmail.com>> .

## Packets

 Each packet is sent over a single socket from the client to the server. A packet consists of one of two things: a fixed length header followed by a variable-length payload. The fixed length header contains all fields that are constant for each packet type, and the payload is defined in the following section under each specific packet type.

 All packets start with the character sequence 33333333333333