<script setup lang="js">
import { ref } from 'vue';

const selectedCurrency = ref({ name: 'EUR', code: 'de-DE' });
const currencies = ref([
    { name: 'EUR', code: 'de-DE' },
    { name: 'USD', code: 'en-US' },
    { name: 'TRY', code: 'tr-TR' },
    { name: 'XMR', code: 'XMR' }
]);
const amount = ref(0);
const amountString = ref(''); // Track the amount as a string for precise cent handling

// Function to handle button click and add digit
function addDigit(digit) {
    // Limit maximum number length for practicality (e.g., 10 digits)
    if (amountString.value.length >= 10) return;

    // Append digit to the amount string
    amountString.value += digit;

    // Convert the string to a float with cents precision
    const numericAmount = parseFloat(amountString.value) / 100;

    // Update the display amount
    amount.value = numericAmount;
}

// Function to clear the amount
function clearAmount() {
    amountString.value = '';
    amount.value = 0;
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
                            <!-- Conditional currency display for XMR or others -->
                            <InputNumber v-if="selectedCurrency.name != 'XMR'" v-model="amount" mode="currency" :currency="selectedCurrency.name" :locale="selectedCurrency.code" size="large" fluid></InputNumber>
                            <InputNumber v-else v-model="amount" suffix=" ɱ" size="large" :minFractionDigits="4" fluid></InputNumber>

                            <Select size="large" optionLabel="name" v-model="selectedCurrency" :options="currencies"></Select>
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
                            <Button size="large" severity="danger" label="CLR" @click="clearAmount"></Button>
                            <Button size="large" label="0" @click="addDigit(0)"></Button>
                            <Button size="large" severity="success" label="ENTER" @click="enterAmount"></Button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
