<script setup lang="js">
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const toast = useToast();
const { push } = useRouter();
const instanceHealth = ref(true);
const instanceVersion = ref('v0.0.0');
const moneropayVersion = ref('v0.0.0');
const isLoading = ref(true);

function getInstance() {
    axios
        .get(import.meta.env.VITE_API_BASE + '/admin/instance', { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            instanceHealth.value = response.data.healthy;
            instanceVersion.value = response.data.version;
            moneropayVersion.value = response.data.moneropay;
        })
        .catch((error) => {
            console.log(error);
            if (error.response) {
                toast.add({ severity: 'error', summary: 'Failed to fetch instance info', detail: error.response.data.message, life: 3000 });
            } else {
                toast.add({ severity: 'error', summary: error.message, detail: error.code, life: 3000 });
            }
        })
        .finally(() => (isLoading.value = false));
}
</script>

<template>
    <div>
        <div class="card">
            <h1 class="font-semibold text-xl mb-5">Account Settings</h1>
            <p><b>Healthy:</b> {{ instanceHealth }}</p>
            <p><b>Version:</b> {{ instanceVersion }}</p>
            <p><b>MoneroPay version:</b> {{ moneropayVersion }}</p>
        </div>
    </div>
</template>
