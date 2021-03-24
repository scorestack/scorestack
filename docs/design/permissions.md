User Roles
==========

Several roles are added to Scorestack during creation. These roles can be assigned to users to give them access to specific parts of the Scorestack components. This document describes the access given to each role, and explains their intended uses.

`dynamicbeat_reader`
--------------------

This role provides read-only access to the `checkdef*` and `attrib_*` indices. This role is intended to be used by the Dynamicbeat user, and provides Dynamicbeat with the least privilege required for proper operation.

`common`
--------

This role provides read-only access to the `results-all*` and `checks` indices, and provides read access to the `scorestack` space. This allows users to view the generic results and some generic checks, which is required for the overall dashboards to work properly. The `scorestack` space is a customized Kibana space that only includes Scorestack-specific components, reducing the clutter in the Kibana UI. Read-only access is provided to only that space so that users don't have to pick between it and the default space, which has several components included that are not needed for Scorestack.

This role should be used for all Scorestack end-users that interact with Kibana.

`spectator`
-----------

This role provides read-only access to the `results*` indices. This allows users to view team-specific dashboards and admin/group check results. This role should generally only be given to Scorestack administrators or "spectator" users like redteam and whiteteam.

`attribute-admin`
-----------------

This role provides full access to the `attrib_*` indices. This allows users to modify all attributes of all teams. This role should only be given to Scorestack administrators that need to modify administrator attributes, or assist teams with modifying their own attributes.

`check-admin`
-------------

This role provides full access to the `check*` indices. This allows users to create, modify, and delete all check definitions. This role should only be given to Scorestack administrators that are managing checks.

Team Roles
----------

A role is created for each team that gets added, which provides read access to the team results index for the team and read/write access to the team's user attributes index. This allows team users to see detailed check results for their team, view team-specific dashboards for their team, and modify user attributes for their team. This role should only be given to the user associated with a team.