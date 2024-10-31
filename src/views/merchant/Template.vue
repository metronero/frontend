<script setup lang="js">
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useSessionStore } from '../../stores/session';

const sessionStore = useSessionStore();
const toast = useToast();
const { push } = useRouter();

const apiBaseUrl = import.meta.env.VITE_API_BASE;

function resetTemplate() {
    axios
        .delete(import.meta.env.VITE_API_BASE + '/merchant/template', { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            toast.add({ severity: 'success', summary: 'Successfully reset template', detail: 'Now using the default template', life: 3000 });
        })
        .catch((error) => {
            console.log(error);
            if (error.response.status == 401) {
                toast.add({ severity: 'error', summary: 'Not logged in', detail: error.code, life: 3000 });
                push({ path: '/auth/login' });
            } else {
                toast.add({ severity: 'error', summary: error.message, detail: error.code, life: 3000 });
            }
        });
}

const files = ref([]);

async function onFileSelect(event) {
    files.value = event.files; // Store all selected files
    const formData = new FormData();

    // Append each file to the FormData object
    files.value.forEach((file) => {
        formData.append('template[]', file); // "template[]" matches backend field name
    });

    try {
        await axios.post(import.meta.env.VITE_API_BASE + '/merchant/template', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
            withCredentials: true
        });

        toast.add({ severity: 'success', summary: 'File upload successful', life: 3000 });
    } catch (error) {
        toast.add({ severity: 'error', summary: 'Failed to upload files', detail: error, life: 3000 });
    }
}
</script>

<template>
    <div>
        <div class="card">
            <h1 class="font-semibold text-2xl mb-5">Template</h1>
            <p>Manage the template used for invoice pages.</p>
            <div class="my-7">
                <Button icon="pi pi-eye" severity="info" as="a" :href="apiBaseUrl + '/merchant/template'" label="Preview current template" />
            </div>

            <h1 class="font-semibold text-xl mb-5">New Template</h1>
            <p>Upload any assets (JS, CSS, images) and the single template page (index.html). In your index.html, reference all assets using the base path:</p>
            <p>"{{ apiBaseUrl + '/merchantdata/' + sessionStore.accountId }}".</p>
            <p class="mt-2">
                For example: <code>{{ '&lt;img src=&quot;' + apiBaseUrl + '/merchantdata/' + sessionStore.accountId + '/mystore.png&quot;&gt;' }}</code>
            </p>
            <div class="my-7">
                <FileUpload name="template[]" customUpload @select="onFileSelect" @upload="onAdvancedUpload($event)" :multiple="true">
                    <template #empty>
                        <span>Drag and drop files to here to upload.</span>
                    </template>
                </FileUpload>
            </div>
            <h1 class="font-semibold text-xl mb-5">Reset Template</h1>
            <p>Delete your current template file from the server and use the default template.</p>
            <div class="my-7">
                <Button icon="pi pi-eye" class="mr-3" as="a" :href="apiBaseUrl + '/merchant/template/default'" severity="info" label="Preview default theme" />
                <Button icon="pi pi-exclamation-triangle" severity="danger" label="Reset current template" @click="resetTemplate()" />
            </div>
        </div>
    </div>
</template>
