<template>
  <div class="index-page">
    <a-input-search
      v-model:value="searchParams.text"
      placeholder="input search text"
      enter-button="Search"
      size="large"
      @search="onSearch"
    />
    {{ JSON.stringify(searchParams) }}
    <MyDivider />
    <a-tabs v-model:activeKey="activeKey" @change="onTabChange">
      <a-tab-pane key="post" tab="文章">
        <PostList />
      </a-tab-pane>
      <a-tab-pane key="user" tab="用户">
        <UserList />
      </a-tab-pane>
      <a-tab-pane key="picture" tab="图片">
        <PictureList />
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script lang="ts" setup>
// 用户在操作时候，改变url，通过url来改变页面状态，不允许反向修改
import { ref, watchEffect } from "vue";
import PostList from "@/components/PostList.vue";
import UserList from "@/components/UserList.vue";
import PictureList from "@/components/PictureList.vue";
import MyDivider from "@/components/MyDivider.vue";
import { useRoute, useRouter } from "vue-router";
import myAxios from "@/plugins/myAxios";

myAxios.post("");

const router = useRouter();
const route = useRoute();
const initSearchParams = ref({
  text: "",
  pageSize: 10,
  pageNum: 10,
});
const searchParams = ref(initSearchParams);
watchEffect(() => {
  searchParams.value = {
    ...initSearchParams,
    text: route.query.text,
  } as any;
});
const activeKey = route.params.category;
const onSearch = (value: string) => {
  alert(value);
  router.push({
    query: searchParams.value,
  });
};

const onTabChange = (key: string) => {
  console.log(key);
  router.push({
    path: `/${key}`,
    query: searchParams.value,
  });
};
</script>
