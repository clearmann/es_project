<template>
  用户列表
  <a-list item-layout="horizontal" :data-source="props.postList">
    <template #renderItem="{ item }">
      <a-card hoverable style="width: 240px">
        <template #cover>
          <img
            alt="example"
            src="https://os.alipayobjects.com/rmsportal/QBnOOoLaAfKPirc.png"
          />
        </template>
        <a-card-meta :title="item.title">
          <template #description>{{ item.content }}</template>
        </a-card-meta>
      </a-card>
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
