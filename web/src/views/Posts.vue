<template>
  <div>
    <NewPost />
    <h1>Posts</h1>
    <div v-for="post in posts" :key="post.id">
      <h2>{{ post.title }}</h2>
      <p>{{ post.body }}</p>
    </div>
  </div>
</template>
<script>
import NewPost from "@/components/NewPost.vue";
export default {
  components: {
    NewPost,
  },
  data() {
    return {
      posts: [],
      page: 1,
      limit: 10,
    };
  },
  created() {
    this.getPosts();
  },
  methods: {
    async getPosts() {
      const response = await fetch(
        `https://jsonplaceholder.typicode.com/posts?_page=${this.page}&_limit=${this.limit}`
      );
      const data = await response.json();
      this.posts = data;
    },
  },
};
</script>
