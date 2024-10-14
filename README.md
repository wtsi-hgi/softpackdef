softpackdef
===========

A script to retrieve singularity definition files from SoftPack images, removing SoftPack specific build cache code. Can also convert to a Dockerfile.

Usage
=====

Can be used in one of two ways; you can specify the SoftPack environment manually:

```bash
module load module/for/softpackdef
softpackdef users/me/myEnv-1.0
```

…or you can load the SoftPack module and let the script autodiscover it:

```bash
module load module/for/softpackdef
module load users/me/myEnv/1.0
softpackdef
```

You can also specify the `--docker` flag to produce a Dockerfile.
