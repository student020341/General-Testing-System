# Testing System

Project idea: a multi user crud application where users define tests as a series of equations. The equations are javascript closures that can have their inputs fed from other equation outputs from the local test or other external tests. A report is a collection of tests.

## Calculation

A calulcation may be declared like any valid javascript closure:
```
// simple
(a, b) => a + b

// with brackets
(first, second) => {
    const third = first + second;
    return third;
}
```

The system should extract variable names from the closure parameters and present them as a list where users can choose how the value is populated, being fed by an input field, output from another calculation in the same test, or a calculation from a different test under the same report.

## project architecture

The project is an exercise of domain driven design, clean code, dependency injection, and nats events. Consult https://github.com/uber-go/guide/blob/master/style.md for style.

## tech stack

- NATS for events
- pocket base or surreal for database, undecided
