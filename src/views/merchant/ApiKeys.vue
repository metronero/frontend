<script setup lang="js">
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const apiKeys = ref([]); // Ensure apiKeys is always an array
const deleteApiKeysDialog = ref(false);
const showApiKeyDialog = ref(false);
const newApiKey = ref(''); // Holds the one-time display key
const selectedApiKeys = ref([]); // Holds selected keys for deletion
const toast = useToast();
const { push } = useRouter();

onMounted(() => {
    fetchApiKeys();
});

const fetchApiKeys = async () => {
    try {
        const response = await axios.get(import.meta.env.VITE_API_BASE + '/merchant/apikey', { withCredentials: true });
        apiKeys.value = response.data || []; // Ensure the response data is set as an array
    } catch (error) {
        if (error.response && error.response.status === 401) {
            toast.add({ severity: 'error', summary: 'Not logged in', detail: error.code, life: 3000 });
            push({ path: '/auth/login' });
        } else {
            toast.add({ severity: 'error', summary: error.message, detail: error.code, life: 3000 });
        }
    }
};

const createApiKey = async () => {
    try {
        const response = await axios.post(import.meta.env.VITE_API_BASE + '/merchant/apikey', {}, { withCredentials: true });

        newApiKey.value = response.data?.key || '';
        const newKeyEntry = { key_id: response.data.key_id, expiry: response.data.expiry };

        // Ensure apiKeys is an array before pushing
        if (!apiKeys.value) apiKeys.value = [];
        apiKeys.value.push(newKeyEntry);

        showApiKeyDialog.value = true;
        toast.add({ severity: 'success', summary: 'API Key Created', detail: 'New API key generated', life: 3000 });
    } catch (error) {
        toast.add({ severity: 'error', summary: 'API Key Creation Failed', detail: error.response?.data?.message || error.message, life: 3000 });
    }
};

const confirmDeleteApiKeys = () => {
    deleteApiKeysDialog.value = true;
};

const deleteSelectedApiKeys = async () => {
    try {
        const deletePromises = selectedApiKeys.value.map((apiKey) => axios.delete(`${import.meta.env.VITE_API_BASE}/merchant/apikey/${apiKey.key_id}`, { withCredentials: true }));
        await Promise.all(deletePromises);

        apiKeys.value = apiKeys.value.filter((apiKey) => !selectedApiKeys.value.includes(apiKey));
        selectedApiKeys.value = [];
        deleteApiKeysDialog.value = false;

        toast.add({ severity: 'success', summary: 'API Keys Deleted', detail: 'Selected API keys deleted', life: 3000 });
    } catch (error) {
        toast.add({ severity: 'error', summary: 'Deletion Failed', detail: error.response?.data?.message || error.message, life: 3000 });
    }
};

const copyToClipboard = (text) => {
    navigator.clipboard
        .writeText(text)
        .then(() => {
            toast.add({ severity: 'success', summary: 'Copied', detail: 'API Key copied to clipboard', life: 3000 });
        })
        .catch((error) => {
            toast.add({ severity: 'error', summary: 'Copy Failed', detail: error.message, life: 3000 });
        });
};
</script>

<template>
    <div>
        <div class="card">
            <h1 class="font-semibold text-xl mb-5">API Keys</h1>

            <Toolbar class="mb-6">
                <template #start>
                    <Button label="New API Key" icon="pi pi-plus" severity="secondary" class="mr-2" @click="createApiKey" />
                    <Button label="Delete" icon="pi pi-trash" severity="secondary" @click="confirmDeleteApiKeys" :disabled="!selectedApiKeys || !selectedApiKeys.length" />
                </template>
            </Toolbar>

            <DataTable :value="apiKeys" v-model:selection="selectedApiKeys" dataKey="key_id" :paginator="true" :rows="10" selectionMode="multiple">
                <Column selectionMode="multiple" style="width: 3rem" :exportable="false"></Column>
                <Column field="key_id" header="Key ID" sortable style="min-width: 16rem">
                    <template #body="slotProps"> {{ slotProps.data.key_id }} </template>
                </Column>
                <Column field="expiry" header="Expiry Date" sortable style="min-width: 12rem">
                    <template #body="slotProps">
                        {{ new Date(slotProps.data.expiry).toLocaleDateString() }}
                    </template>
                </Column>
            </DataTable>
        </div>

        <!-- Dialog for showing the new API Key -->
        <Dialog v-model:visible="showApiKeyDialog" :style="{ width: '450px' }" header="New API Key" :modal="true">
            <div class="flex flex-col gap-4 items-center">
                <p class="text-center">This is your new API Key. Please copy and store it in a secure location as it won't be shown again.</p>
                <div class="flex items-center gap-4">
                    <InputText readonly :value="newApiKey" style="width: 100%" />
                    <Button icon="pi pi-copy" label="Copy" @click="copyToClipboard(newApiKey)" />
                </div>
            </div>
            <template #footer>
                <Button label="Close" icon="pi pi-check" @click="showApiKeyDialog = false" />
            </template>
        </Dialog>

        <!-- Dialog for confirming deletion of selected API Keys -->
        <Dialog v-model:visible="deleteApiKeysDialog" :style="{ width: '450px' }" header="Confirm Deletion" :modal="true">
            <div class="flex items-center gap-4">
                <i class="pi pi-exclamation-triangle !text-3xl" />
                <span>Are you sure you want to delete the selected API keys?</span>
            </div>
            <template #footer>
                <Button label="No" icon="pi pi-times" text @click="deleteApiKeysDialog = false" />
                <Button label="Yes" icon="pi pi-check" @click="deleteSelectedApiKeys" />
            </template>
        </Dialog>
    </div>
</template>
