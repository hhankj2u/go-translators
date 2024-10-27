<template>
  <div id="app">
    <h1>Dictionary Search</h1>
    <input v-model="term" placeholder="Enter word"/>
    <button @click="search">Search</button>

    <!-- Render resultHtml in an iframe for full HTML loading -->
    <iframe :srcdoc="formattedResultHtml" width="100%" height="100%"></iframe>
  </div>
</template>

<script>
import { SearchDictionary } from '../wailsjs/go/main/App';

export default {
  data() {
    return {
      term: 'banana',
      resultHtml: '',
      baseUrl: 'https://dictionary.cambridge.org' // Replace with your actual base URL
    };
  },
  computed: {
    formattedResultHtml() {
      // Inject a base URL tag if it's not already present in the HTML
      if (!this.resultHtml.includes('<base')) {
        return `<base href="${this.baseUrl}">` + this.resultHtml;
      }
      return this.resultHtml;
    }
  },
  methods: {
    async search() {
      try {
        this.resultHtml = await SearchDictionary(this.term);
      } catch (error) {
        console.error("Error fetching dictionary result:", error);
      }
    }
  }
};
</script>

<style scoped>
iframe {
  border: none;
}
</style>
