import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useSessionStore = defineStore('session', () => {
    const role = ref('merchant');
    const name = ref('');
    const accountId = ref('');

    if (localStorage.getItem('account_role')) {
        role.value = localStorage.getItem('account_role');
    }
    if (localStorage.getItem('account_name')) {
        name.value = localStorage.getItem('account_name');
    }
    if (localStorage.getItem('account_id')) {
        accountId.value = localStorage.getItem('account_id');
    }

    function setName(n) {
        name.value = n;
    }

    function setRole(r) {
        role.value = r;
    }

    function setAccountId(a) {
        accountId.value = a;
    }

    return { role, name, accountId, setName, setRole, setAccountId };
});
