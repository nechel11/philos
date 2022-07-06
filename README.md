# Dining philosophers problem
## Problem statement
Five philosophers dine together at the same table. Each philosopher has their own place at the table. Their philosophical problem in this instance is that the dish served is a kind of spaghetti which has to be eaten with two forks.

There is a fork between each plate. Each philosopher can only alternately think and eat. Moreover, a philosopher can only eat their spaghetti when they have both a left and right fork. Thus two forks will only be available when their two nearest neighbors are thinking, not eating. After an individual philosopher finishes eating, they will put down both forks. The problem is how to design a regimen (a concurrent algorithm) such that no philosopher will starve; i.e., each can forever continue to alternate between eating and thinking, assuming that no philosopher can know when others may want to eat or think (an issue of incomplete information). 

## How to run
```
git clone https://github.com/nechel11/philos.git philos
cd philos
go install "golang.org/x/net/websocket" (if you don't have websocket module)
go run main.go
```
## How it works

