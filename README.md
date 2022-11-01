# Timing attack to discover password

Code is written using golang. To execute.

```
go build
./timing-attack >yourlogfile 2>&1
```

## How we are trying to guess correct password:

Our alphabet consists of digits only: 0123456789

Based on the entropy we need to find password with the length between 11 and 14 characters.

We assume that better guess will have longer processing time on the server. Obstacle: noize, caused by network.

So for example, if first 3 digits 412 are correct (position 3), request to check that password will take a little bit longer than 413 or 410.

To check response times for every possible digit from the alphabet in the particular position we spawn 10 parallel POST requests and use average duration (to reduce noize) as a final number.

## Results

So far no correct password has been found... :(