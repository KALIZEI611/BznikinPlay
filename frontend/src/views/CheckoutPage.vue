<template>
  <div class="checkout-page">
    <div class="checkout-container">
      <div class="checkout-steps">
        <div class="step" :class="{ active: currentStep >= 1 }">
          <div class="step-number">1</div>
          <div class="step-label">Выбор консоли</div>
        </div>
        <div class="step" :class="{ active: currentStep >= 2 }">
          <div class="step-number">2</div>
          <div class="step-label">Данные аренды</div>
        </div>
        <div class="step" :class="{ active: currentStep >= 3 }">
          <div class="step-number">3</div>
          <div class="step-label">Подтверждение</div>
        </div>
      </div>

      <div v-if="currentStep === 1 && !preselectedConsole" class="step-content">
        <h2>Выберите консоль для аренды</h2>
        <div class="consoles-grid">
          <div
            v-for="console in consoles"
            :key="console.id"
            class="console-select-card"
            :class="{ selected: selectedConsole?.id === console.id }"
            @click="selectConsole(console)"
          >
            <img
              :src="console.image_url"
              :alt="console.model"
              @error="handleImageError"
            />
            <div class="console-select-info">
              <h3>{{ console.model }}</h3>
              <p class="price">{{ console.price_per_day }} ₽/сутки</p>
              <div
                class="availability"
                :class="{ available: console.is_available }"
              >
                {{ console.is_available ? "✅ Доступна" : "❌ Арендована" }}
              </div>
            </div>
          </div>
        </div>
        <button class="btn-next" @click="nextStep" :disabled="!selectedConsole">
          Далее →
        </button>
      </div>

      <div v-if="preselectedConsole && currentStep === 1" class="step-content">
        <div class="preselected-info">
          <h2>Вы выбрали:</h2>
          <div class="selected-console-display">
            <img
              :src="preselectedConsole.image_url"
              :alt="preselectedConsole.model"
            />
            <div class="selected-info">
              <h3>{{ preselectedConsole.model }}</h3>
              <p class="price">
                {{ preselectedConsole.price_per_day }} ₽/сутки
              </p>
              <div class="availability available">✅ Доступна</div>
            </div>
          </div>
          <button class="btn-next" @click="nextStep">
            Продолжить оформление →
          </button>
        </div>
      </div>

      <div v-if="currentStep === 2" class="step-content">
        <h2>Данные аренды</h2>

        <div class="selected-console-info">
          <img :src="selectedConsole?.image_url" alt="" />
          <div>
            <h3>{{ selectedConsole?.model }}</h3>
            <p>{{ selectedConsole?.price_per_day }} ₽/сутки</p>
          </div>
        </div>

        <form @submit.prevent="nextStep">
          <div class="form-group">
            <label>📅 Дата начала аренды</label>
            <input
              type="date"
              v-model="rentalForm.start_date"
              :min="minDate"
              required
            />
          </div>

          <div class="form-group">
            <label>📅 Дата окончания аренды</label>
            <input
              type="date"
              v-model="rentalForm.end_date"
              :min="rentalForm.start_date"
              required
            />
          </div>

          <div class="form-group">
            <label>📍 Адрес доставки</label>
            <textarea
              v-model="rentalForm.delivery_address"
              placeholder="Укажите полный адрес доставки (город, улица, дом, квартира)"
              rows="3"
              required
            ></textarea>
          </div>

          <div class="form-group">
            <label>📞 Контактный телефон</label>
            <input
              type="tel"
              v-model="rentalForm.phone"
              placeholder="+7 (XXX) XXX-XX-XX"
              required
            />
          </div>

          <div class="form-group">
            <label>💬 Комментарий к заказу (необязательно)</label>
            <textarea
              v-model="rentalForm.comment"
              placeholder="Укажите дополнительную информацию, например, домофон, этаж и т.д."
              rows="2"
            ></textarea>
          </div>

          <div class="price-preview">
            <h4>Предварительный расчет:</h4>
            <div class="price-details">
              <span
                >{{ getDaysCount() }} дней ×
                {{ selectedConsole?.price_per_day }} ₽</span
              >
              <strong>= {{ totalPrice }} ₽</strong>
            </div>
          </div>

          <div class="form-actions">
            <button type="button" class="btn-back" @click="prevStep">
              ← Назад
            </button>
            <button type="submit" class="btn-next">Далее →</button>
          </div>
        </form>
      </div>

      <div v-if="currentStep === 3" class="step-content">
        <h2>Подтверждение заказа</h2>

        <div class="order-summary">
          <h3>Детали заказа</h3>

          <div class="summary-item">
            <span>Консоль:</span>
            <strong>{{ selectedConsole?.model }}</strong>
          </div>

          <div class="summary-item">
            <span>Период аренды:</span>
            <strong
              >{{ formatDate(rentalForm.start_date) }} -
              {{ formatDate(rentalForm.end_date) }}</strong
            >
          </div>

          <div class="summary-item">
            <span>Количество дней:</span>
            <strong>{{ getDaysCount() }} дня(ей)</strong>
          </div>

          <div class="summary-item">
            <span>Стоимость аренды:</span>
            <strong class="total-price">{{ totalPrice }} ₽</strong>
          </div>

          <div class="summary-item">
            <span>Адрес доставки:</span>
            <strong>{{ rentalForm.delivery_address }}</strong>
          </div>

          <div class="summary-item">
            <span>Контактный телефон:</span>
            <strong>{{ rentalForm.phone }}</strong>
          </div>

          <div v-if="rentalForm.comment" class="summary-item">
            <span>Комментарий:</span>
            <strong>{{ rentalForm.comment }}</strong>
          </div>
        </div>

        <div class="delivery-info">
          <h4>🚚 Информация о доставке</h4>
          <p>Доставка осуществляется в день начала аренды с 10:00 до 20:00</p>
          <p>
            Курьер свяжется с вами за час до прибытия по указанному телефону
          </p>
          <p>
            При возврате консоль нужно будет передать курьеру в последний день
            аренды
          </p>
        </div>

        <div class="form-actions">
          <button type="button" class="btn-back" @click="prevStep">
            ← Назад
          </button>
          <button type="button" class="btn-submit" @click="confirmOrder">
            Подтвердить заказ
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";

export default {
  name: "CheckoutPage",
  setup() {
    const router = useRouter();
    const route = useRoute();
    const consoles = ref([]);
    const selectedConsole = ref(null);
    const preselectedConsole = ref(null);
    const currentStep = ref(1);
    const minDate = new Date().toISOString().split("T")[0];

    const rentalForm = ref({
      start_date: "",
      end_date: "",
      delivery_address: "",
      phone: "",
      comment: "",
    });

    const totalPrice = computed(() => {
      if (
        !selectedConsole.value ||
        !rentalForm.value.start_date ||
        !rentalForm.value.end_date
      ) {
        return 0;
      }
      const days = getDaysCount();
      return days * selectedConsole.value.price_per_day;
    });

    const getDaysCount = () => {
      if (!rentalForm.value.start_date || !rentalForm.value.end_date) {
        return 0;
      }
      const start = new Date(rentalForm.value.start_date);
      const end = new Date(rentalForm.value.end_date);
      const diffTime = Math.abs(end - start);
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24)) + 1;
      return diffDays;
    };

    const fetchConsoles = async () => {
      try {
        const response = await axios.get("/api/consoles");
        consoles.value = response.data;

        const consoleId = route.params.id;
        if (consoleId) {
          const foundConsole = consoles.value.find(
            (c) => c.id === parseInt(consoleId) && c.is_available,
          );
          if (foundConsole) {
            preselectedConsole.value = foundConsole;
            selectedConsole.value = foundConsole;
          } else {
            alert("Выбранная консоль недоступна для аренды");
            router.push("/");
          }
        }
      } catch (error) {
        console.error("Error fetching consoles:", error);
      }
    };

    const selectConsole = (console) => {
      selectedConsole.value = console;
    };

    const nextStep = () => {
      if (currentStep.value === 1 && !selectedConsole.value) {
        alert("Пожалуйста, выберите консоль");
        return;
      }

      if (currentStep.value === 2) {
        if (!rentalForm.value.start_date || !rentalForm.value.end_date) {
          alert("Пожалуйста, выберите даты аренды");
          return;
        }
        if (!rentalForm.value.delivery_address.trim()) {
          alert("Пожалуйста, укажите адрес доставки");
          return;
        }
        if (!rentalForm.value.phone.trim()) {
          alert("Пожалуйста, укажите контактный телефон");
          return;
        }
      }

      if (currentStep.value < 3) {
        currentStep.value++;
      }
    };

    const prevStep = () => {
      if (currentStep.value > 1) {
        currentStep.value--;
      }
    };

    const confirmOrder = async () => {
      try {
        const rentalData = {
          console_id: selectedConsole.value.id,
          start_date: new Date(rentalForm.value.start_date).toISOString(),
          end_date: new Date(rentalForm.value.end_date).toISOString(),
          delivery_address: rentalForm.value.delivery_address,
          phone: rentalForm.value.phone,
          comment: rentalForm.value.comment,
        };

        await axios.post("/api/rentals", rentalData);
        alert(
          "Заказ успешно оформлен! Консоль будет доставлена по указанному адресу 🎮\nКурьер свяжется с вами за час до прибытия.",
        );
        router.push("/rentals");
      } catch (error) {
        console.error("Order error:", error);
        alert(error.response?.data?.error || "Ошибка при оформлении заказа");
      }
    };

    const formatDate = (date) => {
      if (!date) return "";
      return new Date(date).toLocaleDateString("ru-RU", {
        day: "numeric",
        month: "long",
        year: "numeric",
      });
    };

    const handleImageError = (event) => {
      event.target.src = "https://via.placeholder.com/200x150?text=Console";
    };

    const setDefaultDates = () => {
      const today = new Date();
      const tomorrow = new Date(today);
      tomorrow.setDate(tomorrow.getDate() + 1);

      rentalForm.value.start_date = tomorrow.toISOString().split("T")[0];
      rentalForm.value.end_date = tomorrow.toISOString().split("T")[0];
    };

    onMounted(() => {
      fetchConsoles();
      setDefaultDates();
    });

    return {
      consoles,
      selectedConsole,
      preselectedConsole,
      currentStep,
      rentalForm,
      totalPrice,
      minDate,
      selectConsole,
      nextStep,
      prevStep,
      confirmOrder,
      getDaysCount,
      formatDate,
      handleImageError,
    };
  },
};
</script>

<style scoped>
.checkout-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8f0fe 100%);
  padding: 2rem;
}

.checkout-container {
  max-width: 900px;
  margin: 0 auto;
  background: white;
  border-radius: 20px;
  padding: 2rem;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.checkout-steps {
  display: flex;
  justify-content: space-between;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #e2e8f0;
}

.step {
  flex: 1;
  text-align: center;
  position: relative;
}

.step-number {
  width: 40px;
  height: 40px;
  background: #e2e8f0;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  margin-bottom: 0.5rem;
  transition: all 0.3s;
}

.step.active .step-number {
  background: #0066cc;
  color: white;
}

.step-label {
  font-size: 0.9rem;
  color: #718096;
}

.step.active .step-label {
  color: #0066cc;
  font-weight: 500;
}

.step-content h2 {
  margin-bottom: 1.5rem;
  color: #2d3748;
}

.consoles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.console-select-card {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 15px;
  cursor: pointer;
  transition: all 0.3s;
}

.console-select-card:hover {
  border-color: #0066cc;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 102, 204, 0.1);
}

.console-select-card.selected {
  border-color: #0066cc;
  background: #f0f7ff;
}

.console-select-card img {
  width: 80px;
  height: 80px;
  object-fit: cover;
  border-radius: 10px;
}

.console-select-info {
  flex: 1;
}

.console-select-info h3 {
  font-size: 1rem;
  margin-bottom: 0.5rem;
  color: #2d3748;
}

.price {
  color: #0066cc;
  font-weight: bold;
  margin: 0.25rem 0;
}

.availability {
  font-size: 0.85rem;
  margin-top: 0.25rem;
}

.availability.available {
  color: #38a169;
}

.preselected-info {
  text-align: center;
}

.selected-console-display {
  display: flex;
  gap: 2rem;
  padding: 2rem;
  background: #f7fafc;
  border-radius: 15px;
  margin: 2rem 0;
  align-items: center;
  justify-content: center;
}

.selected-console-display img {
  width: 150px;
  height: 150px;
  object-fit: cover;
  border-radius: 15px;
}

.selected-info {
  text-align: left;
}

.selected-info h3 {
  font-size: 1.3rem;
  margin-bottom: 0.5rem;
  color: #2d3748;
}

.selected-console-info {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background: #f7fafc;
  border-radius: 15px;
  margin-bottom: 1.5rem;
}

.selected-console-info img {
  width: 100px;
  height: 100px;
  object-fit: cover;
  border-radius: 10px;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #4a5568;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e2e8f0;
  border-radius: 10px;
  font-size: 1rem;
  transition: all 0.3s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #0066cc;
  box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
}

.price-preview {
  background: linear-gradient(135deg, #f0f7ff 0%, #e8f0fe 100%);
  padding: 1rem;
  border-radius: 10px;
  margin: 1rem 0;
}

.price-details {
  display: flex;
  justify-content: space-between;
  font-size: 1.1rem;
}

.price-details strong {
  font-size: 1.3rem;
  color: #0066cc;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}

.btn-back,
.btn-next,
.btn-submit {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-back {
  background: #e2e8f0;
  color: #4a5568;
}

.btn-back:hover {
  background: #cbd5e0;
}

.btn-next,
.btn-submit {
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  color: white;
  flex: 1;
}

.btn-next:hover:not(:disabled),
.btn-submit:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 102, 204, 0.3);
}

.btn-next:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.order-summary {
  background: #f7fafc;
  padding: 1.5rem;
  border-radius: 15px;
  margin-bottom: 1.5rem;
}

.summary-item {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem 0;
  border-bottom: 1px solid #e2e8f0;
}

.summary-item:last-child {
  border-bottom: none;
}

.total-price {
  font-size: 1.2rem;
  color: #0066cc;
}

.delivery-info {
  background: #fff5e6;
  padding: 1rem;
  border-radius: 10px;
  margin-bottom: 1.5rem;
}

.delivery-info h4 {
  margin-bottom: 0.5rem;
  color: #ed8936;
}

.delivery-info p {
  font-size: 0.9rem;
  color: #4a5568;
  margin: 0.25rem 0;
}

@media (max-width: 768px) {
  .checkout-page {
    padding: 1rem;
  }

  .checkout-container {
    padding: 1rem;
  }

  .checkout-steps {
    margin-bottom: 1rem;
  }

  .step-number {
    width: 35px;
    height: 35px;
    font-size: 0.9rem;
  }

  .step-label {
    font-size: 0.75rem;
  }

  .consoles-grid {
    grid-template-columns: 1fr;
  }

  .console-select-card {
    padding: 0.75rem;
  }

  .console-select-card img {
    width: 60px;
    height: 60px;
  }

  .form-actions {
    flex-direction: column;
    gap: 0.75rem;
  }

  .btn-back,
  .btn-next,
  .btn-submit {
    width: 100%;
  }

  .selected-console-display {
    flex-direction: column;
    padding: 1rem;
  }

  .selected-console-display img {
    width: 120px;
    height: 120px;
  }

  .selected-info {
    text-align: center;
  }

  .order-summary {
    padding: 1rem;
  }

  .summary-item {
    flex-direction: column;
    gap: 0.5rem;
    text-align: center;
  }

  .delivery-info p {
    font-size: 0.85rem;
  }
}

@media (max-width: 480px) {
  .checkout-steps {
    flex-direction: column;
    gap: 0.5rem;
  }

  .step {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .step-number {
    margin-bottom: 0;
  }

  .price-details {
    flex-direction: column;
    text-align: center;
    gap: 0.5rem;
  }

  .form-group input,
  .form-group textarea {
    font-size: 0.9rem;
  }
}
</style>
