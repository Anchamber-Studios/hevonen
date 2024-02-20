# Club Service
Implements the API for managing clubs. Clubs have members and basic information. Members can be added to a club, removed from a club, or updated in a club.

## Domain Entities
- **Club**: The entity that represents a club. It represents a horse riding club.
- **Contact**: A person who is in contact with the club. It can be a member, employee  or visitor of the club. A contact can login and manage it's information.
- **Role**: Roles are used to manage the authorizations a contact has within the club. Custom roles can be created but there are already some default roles. The default roles are:
  - **admin**: Can manage the basic information of the club. Can assign roles to other contacts and remove/add new contacts.
  - **manager**: Can manage the basic information and contacts of the club.
  - **trainer**: Can manage the contacts.
  - **member**: Can manage it's own information
  - **guest**: Can manage it's own information

## Actions

### Create a new club
When creating a new club, the user which creates the club is assigned the admin role for the club. To create a club, the following inforations are required at least:
- name of the club
- contact information
  - mail
  - phone number

### Update club information

### Add a member to a club
To create a member at least a first name or email is required. If an email is provided, the system will also create an user for the member.
If the email is alredy used by a user, the user will get a notification.

### Remove a memnber from a club

### Assign roles to a member
There needs to be at least one admin at all times for each club. An admin can assign other members with different roles.
