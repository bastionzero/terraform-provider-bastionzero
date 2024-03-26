## v0.3.1 (March 26, 2024)

NOTES:


* Update `bastionzero_db_target` resource documentation ([#72](https://github.com/bastionzero/terraform-provider-bastionzero/issues/72)).


## v0.3.0 (November 15, 2023)

FEATURES:


* data-source/supported_database_configs: Get a list of all supported database authentication configs ([#67](https://github.com/bastionzero/terraform-provider-bastionzero/issues/67)).


* resource/db_target: Add support for managing Db targets ([#67](https://github.com/bastionzero/terraform-provider-bastionzero/issues/67)).


ENHANCEMENTS:


* data-source/db_target, data-source/db_targets, data-source/web_target, data-source/web_targets: Add `proxy_environment_id` attribute ([#68](https://github.com/bastionzero/terraform-provider-bastionzero/issues/68)).


NOTES:


* Update to Go 1.20 and require it to build and run this provider; this is the last Go release that runs on macOS 10.13 High Sierra, macOS 10.14 Mojave, Windows 7, Windows 8, Windows Server 2008, and Windows Server 2012. A future release of this provider will update to Go 1.21, and these platforms will no longer be supported ([#62](https://github.com/bastionzero/terraform-provider-bastionzero/issues/62)).


* Upgraded [`bastionzero-sdk-go`](https://github.com/bastionzero/bastionzero-sdk-go) to v0.10.0 ([#69](https://github.com/bastionzero/terraform-provider-bastionzero/issues/69)).


## v0.2.0 (October 18, 2023)

ENHANCEMENTS:


* data-source/db_target, data-source/db_targets: Add `database_authentication_config` attribute which contains information about the db target's database auth configuration ([#58](https://github.com/bastionzero/terraform-provider-bastionzero/issues/58)).


NOTES:


* data-source/db_target, data-source/db_targets: Deprecate `database_type` and `is_split_cert` attributes; use `database_authentication_config.database` and `database_authentication_config.authentication_type == "SplitCert"` instead ([#58](https://github.com/bastionzero/terraform-provider-bastionzero/issues/58)).


* Upgraded [`bastionzero-sdk-go`](https://github.com/bastionzero/bastionzero-sdk-go) to v0.9.0 ([#60](https://github.com/bastionzero/terraform-provider-bastionzero/issues/60)).


## v0.1.2 (August 30, 2023)

NOTES:


* Upgraded [`bastionzero-sdk-go`](https://github.com/bastionzero/bastionzero-sdk-go) to v0.6.0 ([#43](https://github.com/bastionzero/terraform-provider-bastionzero/issues/43)).


BUG FIXES:


* resource/targetconnect_policy: Fix `verbs` validation to match behavior of the BastionZero API; `RDP` and `SQLServer` are permitted by the remote API ([#36](https://github.com/bastionzero/terraform-provider-bastionzero/issues/36)).


* resource/environment: Remove require replacement behavior when an environment's `name` changes since the BastionZero API permits editing an environment's name after creation ([#35](https://github.com/bastionzero/terraform-provider-bastionzero/issues/35)).


## v0.1.1 (May 24, 2023)

NOTES:


* Upgraded [`bastionzero-sdk-go`](https://github.com/bastionzero/bastionzero-sdk-go) to v0.2.0 ([#17](https://github.com/bastionzero/terraform-provider-bastionzero/issues/17)).


BUG FIXES:


* resource/environment: Fix `offline_cleanup_timeout_hours` validation to match behavior of the BastionZero API ([#15](https://github.com/bastionzero/terraform-provider-bastionzero/issues/15)).


## v0.1.0 (April 19, 2023)

FEATURES:


* resource/environment: Add support for managing environments ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* resource/jit_policy: Add support for managing JIT policies ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* resource/kubernetes_policy: Add support for managing Kubernetes policies ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* resource/proxy_policy: Add support for managing proxy policies ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* resource/sessionrecording_policy: Add support for managing session recording policies ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* resource/targetconnect_policy: Add support for managing target connect policies ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/ad_bash: Add support for getting bash autodiscovery script ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/bzero_target: Add support for getting Bzero target ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/bzero_targets: Add support for getting list of Bzero targets ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/cluster_target: Add support for getting Cluster target ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/cluster_targets: Add support for getting list of Cluster targets ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/dac_target: Add support for getting dynamic access configuration ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/dac_targets: Add support for getting list of dynamic access configurations ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/db_target: Add support for getting Db target ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/db_targets: Add support for getting list of Db targets ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/environment: Add support for getting environment ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/environments: Add support for getting list of environments ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/groups: Add support for getting list of synced IdP groups ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/jit_policies: Add support for getting list of JIT policies ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/jit_policy: Add support for getting JIT policy ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/kubernetes_policies: Add support for getting list of Kubernetes policies ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/kubernetes_policy: Add support for getting Kubernetes policy ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/proxy_policies: Add support for getting list of proxy policies ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/proxy_policy: Add support for getting proxy policy ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/service_account: Add support for getting service account ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/service_accounts: Add support for getting list of service accounts ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/sessionrecording_policies: Add support for getting list of session recording policies ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/sessionrecording_policy: Add support for getting session recording policy ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/targetconnect_policies: Add support for getting list of target connect policies ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/targetconnect_policy: Add support for getting target connect policy ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/user: Add support for getting user ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/users: Add support for getting list of users ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/web_target: Add support for getting Web target ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


* data-source/web_targets: Add support for getting list of Web targets ([#1](https://github.com/bastionzero/terraform-provider-bastionzero/issues/1)).


## v0.1.0-rc.2 (April 13, 2023)

Prerelease (candidate #2) for v0.1.0

## v0.1.0-rc.1 (April 13, 2023)

Prerelease (candidate #1) for v0.1.0
