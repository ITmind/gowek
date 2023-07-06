<script>
  import svelteLogo from './assets/svelte.svg'
  import djangoLogo from './assets/django.svg'
  import Counter from './lib/Counter.svelte'
  import Notes from "./lib/Notes.svelte";
</script>

<main>
  <div>
    <a href="https://www.djangoproject.com" target="_blank" rel="noreferrer">
      <img src={djangoLogo} class="logo" alt="Django Logo" />
    </a>
    <a href="https://svelte.dev" target="_blank" rel="noreferrer">
      <img src={svelteLogo} class="logo svelte" alt="Svelte Logo" />
    </a>
  </div>
  <h1>Django + Svelte</h1>

  <div class="notes">
    <Notes />
  </div>

  <div class="card">
    <Counter />
  </div>

  <p class="read-the-docs">
    Click on the Django and Svelte logos to learn more
  </p>
</main>

<style>
  .logo {
    height: 6em;
    padding: 1.5em;
    will-change: filter;
    transition: filter 300ms;
  }
  .logo:hover {
    filter: drop-shadow(0 0 2em #646cffaa);
  }
  .logo.svelte:hover {
    filter: drop-shadow(0 0 2em #ff3e00aa);
  }
  .read-the-docs {
    color: #888;
  }
</style>
