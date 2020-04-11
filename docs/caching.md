# caching

ghostdog determines what has been built previously and the resulting outcomes by caching. There are two rules
to focus on in regards to caching.

# creating the cache

First, `files` are never cached. `files` only exist to find and validate filepaths exist and a `files` "output"
is the same list of files passed in via the `paths` argument.

Second, `rule`s have two components to caching. First, a rule's cache directory is created by creating a hash.
That hash is created from:

- the file contents of all outputs of all `rule`s a rule depends on (the `sources` argument)
- the rule's commands
- the rule's outputs' filenames

The rule's cache directory is then populated with the rule's outputs (the `outputs` argument). And the cycle continues.

# using the cache

When ghostdog constructs the rule's cache directory ghostdog checks if the cache directory already exists.
If the cache directory exists then ghostdog uses the cached output populated in the cache directory.
Otherwise, ghostdog will execute the rule's commands.

# caching rules without outputs

Some `rule`s create zero outputs like linting. The same caching logic is used for these `rule`s as well.
The only difference is their cache directories are empty, but the existence of the cache directory lets
ghostdog know this rule has been ran previously.

# cache directory structure

By default, ghostdog writes its cache to `$XDG_CACHE_DIR/ghostdog` or `$HOME/.cache/ghostdog`. Within that
directory are a bunch of directories whose names are limited to two characters. These are the first two characters
of the rule's cache directory. Within any of these directories are more directories with the full hash as the
directory name. This is similar to other programs like `git`. This is often done because some file systems
struggle when a large amount of files exist in a directory, so splitting the directories up like this helps
limit how many files/directories directly exist in any directory.
