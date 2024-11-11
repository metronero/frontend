<script setup lang="js">
import { ref } from 'vue';

// Define the currencies with currency symbols in the 'code' field
const selectedCurrency = ref('EUR');
const currencies = ref(['XMR', 'EUR', 'USD']);

const amount = ref(0);
const amountString = ref(''); // Track the amount as a string for precise handling

// Function to handle button click and add digit
function addDigit(digit) {
    // Limit maximum number length for practicality (e.g., 10 digits)
    if (amountString.value.length >= 10) return;

    // Append digit to the amount string
    amountString.value += digit;
    updateAmountValue();
}

// Function to add two zeros
function addDoubleZero() {
    if (amountString.value.length + 2 > 10) return; // Ensure we don't exceed max length
    amountString.value += '00';
    updateAmountValue();
}

// Function to handle decimal point input
function addDecimal() {
    if (!amountString.value.includes('.')) {
        amountString.value += '.';
    }
}

// Function to clear the amount (used by the "times-circle" button)
function clearAmount() {
    amountString.value = '';
    amount.value = 0;
}

// Function to delete the last character (used by the "delete-left" button)
function deleteLastChar() {
    amountString.value = amountString.value.slice(0, -1);
    updateAmountValue();
}

// Function to update the numeric amount based on the amountString
function updateAmountValue() {
    amount.value = parseFloat(amountString.value) || 0;
}

// Function to handle "Enter" action
function enterAmount() {
    console.log('Amount entered:', amount.value);
    // Additional logic for "ENTER" action can be added here
}
</script>

<template>
    <div>
        <div class="flex items-center justify-center">
            <div class="flex flex-col items-center justify-center">
                <div style="border-radius: 56px; padding: 0.3rem; background: linear-gradient(180deg, var(--primary-color) 10%, rgba(33, 150, 243, 0) 30%)">
                    <div class="w-full bg-surface-0 dark:bg-surface-900 py-20 px-8 sm:px-20" style="border-radius: 53px">
                        <div class="mb-5 flex gap-4">
                            <!-- InputText to display amount with selected currency symbol as suffix -->
                            <InputText v-model="amountString" size="large" fluid readonly></InputText>
                            <Select size="large" v-model="selectedCurrency" :options="currencies"></Select>
                        </div>
                        <div class="grid grid-cols-3 gap-4 !text-2xl">
                            <Button size="large" label="1" @click="addDigit(1)"></Button>
                            <Button size="large" label="2" @click="addDigit(2)"></Button>
                            <Button size="large" label="3" @click="addDigit(3)"></Button>
                            <Button size="large" label="4" @click="addDigit(4)"></Button>
                            <Button size="large" label="5" @click="addDigit(5)"></Button>
                            <Button size="large" label="6" @click="addDigit(6)"></Button>
                            <Button size="large" label="7" @click="addDigit(7)"></Button>
                            <Button size="large" label="8" @click="addDigit(8)"></Button>
                            <Button size="large" label="9" @click="addDigit(9)"></Button>
                            <Button size="large" label="00" @click="addDoubleZero"></Button>
                            <Button size="large" label="0" @click="addDigit(0)"></Button>
                            <Button size="large" label="." @click="addDecimal"></Button>
                            <Button size="large" severity="danger" icon="pi pi-times-circle text-2xl" label=" " @click="clearAmount"></Button>
                            <Button size="large" severity="info" icon="pi pi-delete-left text-2xl" label=" " @click="deleteLastChar"></Button>
                            <Button size="large" severity="success" icon="pi pi-arrow-right text-2xl" label=" " @click="enterAmount"></Button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
