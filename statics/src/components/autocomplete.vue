<template>
    <div style="position:relative" 
         v-bind:class="{'open': openSuggestion}">
        <input class="form-control input-sm" type="text"
               :placeholder="placeholder"
               :value="value"
               @keydown.enter.prevent="enter"
               @keydown.tab="complete"
               @keydown.down="down"
               @keydown.up="up"
               @input="change($event.target.value)"
        />
        <ul class="dropdown-menu">
            <li v-for="(suggestion, index) in matches"
                v-bind:class="{'active': isActive(index)}"
                @click="click(index)">
                <a href="#">{{ suggestion }}</a>
            </li>
        </ul>
    </div>
</template>

<script>
export default {

  props: {

    value: {
      type: String,
      required: true,
    },

    suggestions: {
      // can be Array or Function 
      // that returns a promise
      required: true,
    },

    placeholder: {
      type: String,
    },

  },

  data() {
    return {
      open: false,
      current: 0,
      fetchedSuggestions: [],
    }
  },

  computed: {

    openSuggestion() {
      return this.value !== "" &&
             this.matches.length !== 0 &&
             this.open === true;
    },

    matches() {
      return this.fetchedSuggestions.filter(s => {
        return s.indexOf(this.value) >= 0;
      });
    },

  },

  methods: {

    fetchSuggestions() {
      if (this.suggestions instanceof Function) {
        this.suggestions().then(data => {
          this.fetchedSuggestions = data;
        });
      }
      else if (this.suggestions instanceof Array) {
        this.fetchedSuggestions = this.suggestions;
      }
    },

    enter() {
      this.complete();
      this.$emit('select', this.matches[this.current]);
    },

    complete() {
      if (this.openSuggestion === true) {
        var value = this.matches[this.current] || this.value;
        this.open = false;
        this.$emit('input', value);
      }
    },

    click(index) {
      this.open = false;
      this.$emit('input', this.matches[index]);
      this.$emit('select', this.matches[index]);
    },

    up() {
      if (this.current > 0)
        this.current--;
    },

    down() {
      if (this.current < this.matches.length - 1)
        this.current++;
    },

    isActive(index) {
      return index == this.current;
    },

    change(value) {
      // FIXME: debounce
      this.fetchSuggestions();
      if (this.open === false) {
        this.open = true;
        this.current = 0;
      }
      this.$emit('input', value);
    }

  }

}
</script>
