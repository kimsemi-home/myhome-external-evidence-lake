# Bootstrap Checklist

- [x] Repository is public under `kimsemi-home`.
- [x] Upstream SSOT is declared as `kimsemi-home/myhome-jarvis`.
- [x] Context pack version is `v1`.
- [x] Ontology version is `concept-registry/v1`.
- [x] Authority contract is `authority/v1`.
- [x] Security contract is `security/v1`.
- [x] Raw payloads are not present.
- [x] Private local lake archives are not present.
- [x] Public-safety verification script is present.
- [x] GitHub Actions validates generated contract metadata.
- [x] Hash cache inputs are declared per bootstrap unit.

## Next Allowed Work

After this skeleton is merged, future work should add one capability at a time:

1. upstream context pack importer
2. public fixture-only lake schema
3. cache-aware generated artifact validator
4. adapter proposal packet without network writes
5. scoped collector implementation after separate approval
