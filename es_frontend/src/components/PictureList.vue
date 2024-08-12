<template>
  图片列表
  <a-list item-layout="horizontal" :data-source="props.postList">
    <template #renderItem="{ item }">
      <a-list-item>
        <a-list-item-meta :description="item.content">
          <template #title>
            <a href="https://www.antdv.com/">{{ item.title }}</a>
          </template>
          <template #avatar>
            <a-avatar :src="item.url" />
          </template>
        </a-list-item-meta>
      </a-list-item>
    </template>
  </a-list>
</template>

<script lang="ts" setup>
import { withDefaults, defineProps } from "vue";
interface Props {
  postList: any[];
}
const props = withDefaults(defineProps<Props>(), {
  postList: () => [],
});
</script>
<style scoped></style>
