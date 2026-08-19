# viam-system

## cmdsensor
```
{
    "cmd" : "<path to command>",
    "args" : [ "a1", "a2" ]
}
```

## apt

Makes sure a list of debian packages is installed on the machine.
Missing packages are installed in the background at startup (and on reconfigure).
`Readings` reports the installed version of each package, or `missing`.

```
{
    "packages" : [ "curl", "htop" ],
    "update" : true
}
```

- `packages` (required): debian package names to keep installed
- `update` (optional): run `apt-get update` before installing missing packages

`DoCommand` with `{"install": true}` re-runs the install check.
