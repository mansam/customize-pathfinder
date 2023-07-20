# Pathfinder Customization

## Intro

This tool facilitates customizing the Pathfinder questionnaire by taking an input questionnaire described in YAML
and generating a SQL migration. This migration can then be built into a custom Pathfinder image and used to replace
the one deployed by the Konveyor operator. The migration will mark the old questions as deleted while leaving them
in place so that they can continue to be referenced by any assessments that had been made with the old questionnaire.

## Instructions

1. Copy and alter the `questionnaire_template.yaml` file found in the root of this repository,
or write a new one following its format. See "Questionnaire Format" below.
2. Run `make` to build the tool, which will be output to `bin/migration`
3. Invoke the tool with `bin/migration questionnaire.yaml > output.sql`
4. Check out the [Pathfinder repository](github.com/konveyor/tackle-pathfinder).
5. Identify the last migration in the Pathfinder repo, which is located under `src/main/resources/db/migration`. Copy the
`output.sql` to the same directory and rename it to match the naming convention, incrementing the version by one. 
At the time of writing, the latest migration file is named `V0.0.5_001__update_assessment_bulk.sql`, so something like
`V0.0.5_002__replace_questionnaire.sql` would be suitable.
6. Build a new Pathfinder image and push it to your container repository of choice.
7. Edit the Tackle CR in your cluster and add the field `pathfinder_image_fqin:` to the spec, with the value pointing
to the location of the new Pathfinder image. You may need to delete the Pathfinder deployment from the Konveyor namespace
to have the operator pick up the changes.
8. If all has gone well, your questionnaire has been deployed. Changes will be reflected in newly created assessments,
and old assessments will continue to reflect the old questions.

## Questionnaire Format

The input to the migration tool must be a YAML file following the below format. A few notes on the format:

* Categories correspond to pages of questions.
* Question `name` is a short, arbitrary internal identifier that will not be displayed to the user.
* Question `description` is displayed in Konveyor as the question's tooltip.
* Option `risk` must be one of `UNKNOWN`, `RED`, `AMBER`, `GREEN`. If `risk` is not in this set, the migration
may appear to succeed but attempting to create a new assessment will fail with a server error as the value will
not be deserialized correctly. This can be resolved by fixing the value in the questionnaire yaml, generating a new
migration, and deploying an updated image.
* `order` controls the display order within the relevant collection.

```yaml
categories:
  - title: Bridgekeeper Questions
    order: 0
    questions:
      - name: NAMES
        order: 0
        question: What... is your name?
        description: What is the name of the knight being quizzed?
        options:
          - option: Sir Lancelot of Camelot
            order: 0
            risk: GREEN
          - option: Sir Robin of Camelot
            order: 1
            risk: AMBER
          - option: Sir Gallahad of Camelot
            order: 2
            risk: RED
          - option: Arthur, King of the Britons
            order: 3
            risk: UNKNOWN
```

A short example yaml and corresponding example SQL output can be found in the `example` directory.
