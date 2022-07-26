# Contributing

Goonware is comprised of two components, the configurator and the daemon. The configurator is for configuring
(hence configuator hurrr) and the daemon is what does the Goonware-ing.

Execution begins in main.go::main(). If --daemon is provided, the daemon is spawned. Otherwise the configurator
is spawned.

All of the code for the configurator (including types for marshalling packages etc.) is in `configurator/`,
likewise all of the daemon code is in `daemon/`.

# License

Please note that by contributing to Goonware you agree to license your contributions under Unlicense,
i.e. you waive all rights you hold to your contributions to the greatest extent permissible by law. See LICENSE
for more details.