@springboot
Feature: Build process-mongodb-persistence-springboot images

  Background:
    Given Clone Kogito examples into local directory

  Scenario: Build process-mongodb-persistence-springboot image
    Then Local example service "process-mongodb-persistence-springboot" is built by Maven and deployed to runtime registry

  Scenario Outline: Build process-mongodb-persistence-springboot image with profile <profile>
    Then Local example service "process-mongodb-persistence-springboot" is built by Maven and deployed to runtime registry with Maven configuration:
      | profile | <profile> |

    Examples:
      | profile |
      | events  |
