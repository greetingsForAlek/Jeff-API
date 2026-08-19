# The Jeff API

The API That contains all of the characters from Jeff.

## HOW IT WORKS
Like any other API, it is capable of getting the character/s, creating them, updating them and deleting them.

Sections:
- [Searching](#searching-for-characters)
- [Creating](#creating-characters)
- [Updating](#updating-characters)
- [Deleting](#deleting-characters)

## SEARCHING FOR CHARACTERS
Searching for character/s is easy, and there are 2 ways to do it:
- Get all of the characters
- Get a single character

To get All of the characters, you can GET query the `/characters` path. This will return all of the characters.

To get a SINGLE character, you can, again, GET query `/characters`, but this time add another slash and the id of the character you want (See more about ids in [the create section](#creating-characters)). So the URL should look something like `/characters/1`.

There are a few more Parameters that can be used, such as name, alignment or sort.

### Name:
Name should look like: `/characters?name=Jeff`

### Alignment:
Alignment should look like: `/characters?alignment=Neutral`

### Canon:
Canon should look like: `/characters?canon=true`

### Sort
Sort sorts by a specific field, such as:
`/characters?sort=name`, `characters?sort=alignment` or `characters?sort=id`.
Normally it sorts in ascending order, but you can switch it to descending with order:

### Order
Order must be used with sort.
you can use asc (ascending) like this: `/characters?sort=name&order=asc`,
or you can use desc (descending) like this: `characters?sort=name&order=desc`

## CREATING CHARACTERS
To create a character, you can POST query `/characters`. You must send a JSON object containing the name, description, alignment, image and wether the character is canon or not. Alignment is normally Good, Bad or Neutral. The Image is normally a link to an image, which could be just an image from the web or from a site like Cloudinary or sites like that. Canon is just a boolean (true/false).

### The Request
The JSON should look like this:

`{"name":"Jeff", "description":"The main character", "alignment":"Good", "image": "image-url", "canon": true}`.

Obviously, replace image-url with the URL to your image.

## UPDATING CHARACTERS
To update a character, you can do an UPDATE query to `/characters/idhere`. obviously replace idhere with
the id of the character that you want to update. You have to send the request JSON like you do for creating a character.

## DELETING CHARACTERS
To delete a character, it is relativley simple, just DELETE query `/characters/idhere`, again, replacing idhere with the id of the character you want to delete.