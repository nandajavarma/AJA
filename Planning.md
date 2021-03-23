# Activity Journaling App

## Requirements
### Activities 
The AJA should be used to track the activities and the time investment in each. In order to do so, users will have buttons to add a certain activity (or sub-activity) to the "daily bucket".
No particular ordering within the day (or within the morning-afternoon-evening? TBD) is required.
#### Activities hierarchy
Two hierarchical levels of activities should be present (Activity and Sub-Activity). Users should be able to create a new Activity or to create a Sub-Activity by indicating a "parent" activity.
Activities should be color-coded. Sub-Activities could have a slightly variation of the original code (some random RGB added)
#### Activities Annotations
For each activity, for each day, an annotation could be added.

### Daily Buckets
Since this doesn't work as a calendar, each daily bucket does not have a "hourly" division. It could possibly be divided in morning-afternoon-evening.
#### Adding Activities
The user should be able to visualize all the "created" activities and add them on the daily bucket in forms of 1-hour blocks. 
The activities should be selected with one of the following (TBD)
* A scroll-down menu
* A tag cloud, with the the more used activities visible
* A search box, possibly with auto-complete?
#### Activities in Bucket
Each 1-hour block in the bucket could stack with the already present block of the same activity creating an n-hours block. The amount of total hours should be visible, either with a number or through the size of the block (TBD)
#### Annotations on Activities
Each n-hours block can have one annotation on it.

### Visualization
#### Zoom-Out
Users should be able to have a broader view of their activity, through some "monthly view" (Which could be only color-coded, no text).
#### Analytics
Users should be able to visualize absolute number of hours spent in each Activity and proportion with the other activities.
