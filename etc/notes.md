miniplan - software planning barebones

Lack of clear problem statements in an implementation independant
manner has negative effects on software development in general.

The goal of this project is to aid requirement engineers to define
problem statements, elicit requirements and manage them effectively
over the life span of a software solution.


## Problems

Multiple solutions exist to one problem, knowing why a certain
solution was rejected is equally important to why another was choosen.

Solutions(features) disconnected from a problem are difficult to
design correctly.

Solutions without requirements tend to drift away from the original
problem, thus creating new problems.

Analysing solutions and requirements in isolation is not enough,
engineers must also see the whole picture to align and value
requirements against each other.

Solutions may introduce new problems, the problematic sum of these
should not exceed the original problem.

Problem statement change over time, software developers must
understand the change to appropriately adapt current solution.


## Solution 1 - miniplan web service

Web service for editing changes.

### Quick start

    $ go install github.com/gregoryv/cmd/miniplan@latest
    $ miniplan -h
