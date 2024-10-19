import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useSessionStore = defineStore('session', () => {
    const role = ref('merchant');
    const name = ref('');

    if (localStorage.getItem('account_role')) {
        role.value = localStorage.getItem('account_role');
    }
    if (localStorage.getItem('account_name')) {
        name.value = localStorage.getItem('account_name');
    }

    function setName(n) {
        name.value = n;
    }

    function setRole(r) {
        role.value = r;
    }

    return { role, name, setName, setRole };
});
