This is the **Bashy** release. The binaries are burnt using following conventions:
- bashy.OSNAME are general purpose binaries valid for most distribution
- bashy.ARCH.OSNAME are binaries compiled for the couple Archiecture/OS (Es. Linux x86)

There isn't any plan to relese distribution specific installer. Anyway,  GOLANG is quite compatible with most OS and distribution, so Bashy can run quite everywhere as it use only pure GO constructs.

If there isn't any option for you OS you can:
1. compile it by yourself following istruction in the project README
2. Enanche the build process to include your OS in the pipeline. To do this, just fork the repo, change the Taskfile.yaml and make a pullrequest

