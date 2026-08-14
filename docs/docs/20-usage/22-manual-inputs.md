# Manual run inputs

Manual workflows can declare typed inputs that the web UI renders before a pipeline is started.

```yaml
when:
  - event: manual

manual:
  inputs:
    action:
      type: choice
      description: Production action
      options:
        - deploy
        - seed
      default: deploy
      required: true

    force_reconfigure:
      type: boolean
      description: Reconfigure running containers
      default: false
```

The initial implementation supports `string`, `choice`, and `boolean` inputs.

Submitted values use Woodpecker's existing manual pipeline variables mechanism. Workflows without `manual.inputs` keep the existing free-form key/value editor for backward compatibility.
