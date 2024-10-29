<template>
  <div id="app">
    <!-- Container for input and button -->
    <div class="input-container">
      <input v-model="term" placeholder="Enter word" @keyup.enter="search" />
      <button @click="search">Search</button>
      <label>
        <input type="checkbox" v-model="disableClipboard" /> Disable Clipboard
      </label>
    </div>

    <!-- Render resultHtml in an iframe for full HTML loading -->
    <div class="results-container">
      <div v-for="(html, name) in results" :key="name" class="result-item">
        <h2>{{ name }} Dictionary</h2>
        <iframe :srcdoc="formattedResultHtml(html, name)"></iframe>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import { SearchDictionary, ReadClipboard } from '../wailsjs/go/main/App';

export default {
  data() {
    return {
      term: '',
      results: {},
      disableClipboard: ref(true),
      previousClipboardContent: '',
      baseUrls: {
        WEBSTER: 'https://www.merriam-webster.com',
        CAMBRIDGE: 'https://dictionary.cambridge.org',
        SOHA: 'http://tratu.soha.vn'
      }
    };
  },
  computed: {
    formattedResultHtml() {
      return (html, name) => {
        const baseUrl = this.baseUrls[name];
        // Inject a base URL tag if it's not already present in the HTML
        if (!html.includes('<base')) {
          return `<base href="${baseUrl}">` + html;
        }
        return html;
      };
    }
  },
  methods: {
    async search() {
      try {
        this.results = await SearchDictionary(this.term);
      } catch (error) {
        console.error('Error fetching dictionary result:', error);
      }
    },
    async readClipboard() {
      if (!this.disableClipboard) {
        try {
          const text = await ReadClipboard();
          if (text !== this.previousClipboardContent) {
            this.previousClipboardContent = text;
            const words = text.split(/\s+/).slice(0, 3).join(' ');
            console.log('Clipboard contents:', words);
            this.term = words;
            this.search();
          }
        } catch (err) {
          console.log('Failed to read clipboard contents: ', err);
        }
      }
    },
  },
  mounted() {
    setInterval(this.readClipboard, 2000);
  }
};
</script>

<style scoped>
#app {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  /* Center children vertically */
  height: 100vh;
  /* Set full height for the app */
}

.input-container {
  display: flex;
  /* Use flex to align input and button */
  margin-bottom: 20px;
  /* Space between input/button and results */
  justify-content: center;
  /* Center children horizontally */
}

.input-container input {
  margin-right: 10px;
  /* Space between input and button */
  padding: 10px;
  flex: 1;
  /* Allow input to grow and take available space */
}

.results-container {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  width: 100%;
  justify-content: space-evenly;
  height: calc(100vh - 80px);
  /* Adjust for input/button height */
  overflow-y: auto;
  /* Add scrolling if content exceeds height */
}

.result-item {
  flex: 1 1 30%;
  /* Each item takes up roughly a third of the container */
  max-width: 30%;
  display: flex;
  flex-direction: column;
  height: 100%;
  /* Full height of the container */
}

iframe {
  flex: 1;
  border: none;
  height: 100%;
  /* Full height within result-item */
  width: 100%;
}
</style>
