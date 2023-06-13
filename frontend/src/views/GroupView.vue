<template>
  group id transfered using pinia: {{ groupStore.groupId }}
  <!-- todo: implement view "GroupView.vue" -->
  <pre style="text-align: left;">
    one group view gap.
    view must include:
    -------------------------------------------------
    - if NO MEMBERSHIP, then:
    - - button to make request to join group
    - - title + description of the group
    -------------------------------------------------
    - if MEMBERSHIP, then:
    - - title + description of the group
    - - open group chat button -> opens "ChatView.vue"
    - - invite user button -> opens "GroupInviteView.vue"
    - - group posts button -> opens "PostsView.vue" , but with posts only from this group
    
    - - create new even section
    - - - event title - input text
    - - - event date+time - input date + input time(or some combined Vue widget)
    - - - event description - textarea
    - - - radio button with two variants 0 - not going to event, 1 - going to event
    - - - - (for creator too, inside event creation section,
    and default value for creator is going to event)
    - - - submit button to create event

    - - list of group events section
    - - - for each user, every event has:
     title, date+time, description, going to event button.
     To prevent spamming , only once choice are allowed.
     After that, going button is disabled/replaced by text/or restyled to green etc.
     So user can only once choose to go to event.
  </pre>
</template>

<script lang="ts" setup>
import { onMounted } from 'vue';
import { useGroupStore } from '@/store/group';
const groupStore = useGroupStore();
</script>