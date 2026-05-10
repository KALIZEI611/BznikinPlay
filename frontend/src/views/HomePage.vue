<template>
  <div class="home">
    <div class="hero">
      <h1>Аренда игровых консолей</h1>
      <p>PlayStation и Xbox с доставкой на дом</p>
      <div class="hero-badges">
        <span class="badge">🎮 Бесплатная доставка</span>
        <span class="badge">⭐ 24/7 Поддержка</span>
        <span class="badge">🔒 Безопасная оплата</span>
      </div>
    </div>

    <div class="filter-section">
      <div class="filter-buttons">
        <button
          @click="filterType = 'all'"
          :class="['filter-btn', { active: filterType === 'all' }]"
        >
          Все консоли
        </button>
        <button
          @click="filterType = 'PS5'"
          :class="['filter-btn', { active: filterType === 'PS5' }]"
        >
          🎮 PlayStation 5
        </button>
        <button
          @click="filterType = 'PS4'"
          :class="['filter-btn', { active: filterType === 'PS4' }]"
        >
          🎮 PlayStation 4
        </button>
        <button
          @click="filterType = 'XBOX'"
          :class="['filter-btn', { active: filterType === 'XBOX' }]"
        >
          🎯 Xbox
        </button>
      </div>
    </div>

    <div class="consoles-grid">
      <div
        v-for="console in filteredConsoles"
        :key="console.id"
        class="console-card"
      >
        <div class="console-image" @click="showDetails(console)">
          <img
            :src="console.image_url"
            :alt="console.model"
            @error="handleImageError"
            class="console-img"
            loading="lazy"
          />
          <div class="console-type" :class="getConsoleClass(console.type)">
            {{ console.type }}
          </div>
          <div class="quick-view">
            <span>🔍 Быстрый просмотр</span>
          </div>
        </div>
        <div class="console-info">
          <h3 @click="showDetails(console)" class="clickable-title">
            {{ console.model }}
          </h3>
          <p class="price">{{ console.price_per_day }} ₽<span>/день</span></p>
          <p class="description">
            {{ getShortDescription(console) }}
          </p>
          <div
            class="availability"
            :class="{ available: console.is_available }"
          >
            <span class="status-dot"></span>
            {{ console.is_available ? "Доступна" : "Арендована" }}
          </div>
          <div class="button-group">
            <button @click="showDetails(console)" class="details-btn">
              Подробнее
            </button>
            <button
              v-if="console.is_available"
              @click="goToCheckout(console)"
              class="rent-btn"
            >
              Арендовать →
            </button>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="showDetailsModal"
      class="modal"
      @click.self="showDetailsModal = false"
    >
      <div class="modal-content details-modal">
        <button class="modal-close" @click="showDetailsModal = false">✕</button>

        <div class="details-container">
          <div class="details-image">
            <img
              :src="selectedConsole?.image_url"
              :alt="selectedConsole?.model"
            />
            <div
              class="console-badge"
              :class="getConsoleClass(selectedConsole?.type)"
            >
              {{ selectedConsole?.type }}
            </div>
          </div>

          <div class="details-info">
            <h2>{{ selectedConsole?.model }}</h2>
            <div class="details-price">
              <span class="price-label">Цена аренды:</span>
              <span class="price-value"
                >{{ selectedConsole?.price_per_day }} ₽/день</span
              >
            </div>

            <div class="details-section">
              <h3>🎮 Что вы получите:</h3>
              <ul>
                <li
                  v-for="item in getConsoleBenefits(selectedConsole?.type)
                    .items"
                  :key="item"
                >
                  ✅ {{ item }}
                </li>
              </ul>
            </div>

            <div class="details-section">
              <h3>🎯 Популярные игры:</h3>
              <div class="games-list">
                <span
                  v-for="game in getConsoleGames(selectedConsole?.type)"
                  :key="game"
                  class="game-tag"
                >
                  🎲 {{ game }}
                </span>
              </div>
            </div>

            <div class="details-section">
              <h3>✨ Почему стоит арендовать у нас:</h3>
              <div class="features-grid">
                <div class="feature">
                  <span class="feature-icon">🎮</span>
                  <span>Попробуйте перед покупкой</span>
                </div>
                <div class="feature">
                  <span class="feature-icon">👥</span>
                  <span>Идеально для вечеринок с друзьями</span>
                </div>
                <div class="feature">
                  <span class="feature-icon">🚚</span>
                  <span>Бесплатная доставка и настройка</span>
                </div>
                <div class="feature">
                  <span class="feature-icon">🎁</span>
                  <span>Подарочный набор в подарок</span>
                </div>
              </div>
            </div>

            <div class="details-section">
              <h3>💭 Отзывы наших клиентов:</h3>
              <div class="reviews">
                <div class="review">
                  <div class="review-header">
                    <span class="reviewer">Алексей</span>
                    <span class="rating">★★★★★</span>
                  </div>
                  <p>
                    "Арендовал PS5 на выходные с друзьями - это было невероятно!
                    Консоль привезли вовремя, все настроили. Играли в FIFA и
                    Call of Duty до утра. Обязательно возьму еще!"
                  </p>
                </div>
                <div class="review">
                  <div class="review-header">
                    <span class="reviewer">Мария</span>
                    <span class="rating">★★★★★</span>
                  </div>
                  <p>
                    "Заказала Xbox Series X для сына на день рождения. Он был в
                    восторге! Отличное состояние консоли, все игры работают
                    идеально. Спасибо за праздник!"
                  </p>
                </div>
              </div>
            </div>

            <div class="details-actions">
              <button
                v-if="selectedConsole?.is_available"
                @click="goToCheckout(selectedConsole)"
                class="btn-primary rent-now-btn"
              >
                Арендовать сейчас →
              </button>
              <button @click="showDetailsModal = false" class="btn-secondary">
                Закрыть
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch, computed } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";

export default {
  name: "HomePage",
  setup() {
    const router = useRouter();
    const consoles = ref([]);
    const filterType = ref("all");
    const showDetailsModal = ref(false);
    const selectedConsole = ref(null);

    const filteredConsoles = computed(() => {
      if (filterType.value === "all") {
        return consoles.value;
      }
      return consoles.value.filter((c) => c.type === filterType.value);
    });

    const fetchConsoles = async () => {
      try {
        const response = await axios.get("/api/consoles");
        consoles.value = response.data;
      } catch (error) {
        console.error("Error fetching consoles:", error);
      }
    };

    const handleImageError = (event) => {
      event.target.src =
        "https://via.placeholder.com/400x300?text=Console+Image";
    };

    const goToCheckout = (console) => {
      router.push(`/checkout/${console.id}`);
    };

    const getConsoleClass = (type) => {
      if (type === "PS5") return "ps5";
      if (type === "PS4") return "ps4";
      if (type === "XBOX") return "xbox";
      return "";
    };

    const getConsoleBenefits = (type) => {
      const benefits = {
        PS5: {
          items: [
            "Консоль в идеальном состоянии",
            "2 контроллера DualSense",
            "Зарядная станция для контроллеров",
            "Кабель HDMI 2.1",
            "Подписка PlayStation Plus на месяц",
          ],
        },
        PS4: {
          items: [
            "Консоль с новым термоинтерфейсом",
            "1 контроллер DualShock 4",
            "Большой выбор цифровых игр",
            "Кабель HDMI в комплекте",
            "Подборка лучших эксклюзивов",
          ],
        },
        XBOX: {
          items: [
            "Консоль как новая с полной диагностикой",
            "1 контроллер Xbox Wireless",
            "Game Pass Ultimate на месяц",
            "Кабель HDMI 2.1",
            "Возможность быстрого возобновления игр",
          ],
        },
      };
      return benefits[type] || benefits.PS5;
    };

    const getConsoleGames = (type) => {
      const games = {
        PS5: [
          "Spider-Man 2",
          "God of War Ragnarok",
          "Horizon Forbidden West",
          "FIFA 24",
          "Call of Duty MW3",
          "Gran Turismo 7",
        ],
        PS4: [
          "The Last of Us 2",
          "God of War",
          "Red Dead Redemption 2",
          "GTA V",
          "Uncharted 4",
          "Bloodborne",
        ],
        XBOX: [
          "Halo Infinite",
          "Forza Horizon 5",
          "Starfield",
          "Gears 5",
          "Microsoft Flight Simulator",
          "Sea of Thieves",
        ],
      };
      return games[type] || games.PS5;
    };

    const showDetails = (console) => {
      selectedConsole.value = console;
      showDetailsModal.value = true;
    };

    const getShortDescription = (console) => {
      const descriptions = {
        PS5: "Новейшее поколение консолей с потрясающей графикой и скоростью загрузки",
        PS4: "Огромная библиотека игр и доступная цена аренды",
        XBOX: "Мощная консоль с Game Pass и отличной совместимостью",
      };
      return (
        descriptions[console.type] ||
        console.description.substring(0, 80) + "..."
      );
    };

    onMounted(() => {
      fetchConsoles();
    });

    return {
      consoles,
      filteredConsoles,
      filterType,
      showDetailsModal,
      selectedConsole,
      showDetails,
      goToCheckout,
      handleImageError,
      getConsoleClass,
      getConsoleBenefits,
      getConsoleGames,
      getShortDescription,
    };
  },
};
</script>

<style scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.hero {
  text-align: center;
  padding: 4rem 2rem;
  background: linear-gradient(135deg, #f0f7ff 0%, #e8f0fe 100%);
  border-radius: 30px;
  margin-bottom: 3rem;
  border: 1px solid rgba(0, 102, 204, 0.1);
}

.hero h1 {
  font-size: 3rem;
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  margin-bottom: 1rem;
}

.hero p {
  font-size: 1.2rem;
  color: #4a5568;
  margin-bottom: 2rem;
}

.hero-badges {
  display: flex;
  justify-content: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.badge {
  background: white;
  padding: 0.5rem 1rem;
  border-radius: 25px;
  color: #0066cc;
  font-size: 0.9rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.filter-section {
  margin-bottom: 2rem;
}

.filter-buttons {
  display: flex;
  justify-content: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.filter-btn {
  padding: 0.75rem 1.5rem;
  background: white;
  border: 2px solid #e2e8f0;
  border-radius: 25px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  color: #4a5568;
  transition: all 0.3s;
}

.filter-btn:hover {
  border-color: #0066cc;
  color: #0066cc;
}

.filter-btn.active {
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  color: white;
  border-color: transparent;
}

.consoles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 2rem;
}

.console-card {
  background: white;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
  border: 1px solid rgba(0, 102, 204, 0.1);
  display: flex;
  flex-direction: column;
}

.console-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 40px rgba(0, 102, 204, 0.15);
}

.console-image {
  position: relative;
  background: linear-gradient(135deg, #f7fafc 0%, #edf2f7 100%);
  height: 250px;
  width: 100%;
  overflow: hidden;
  flex-shrink: 0;
  cursor: pointer;
}

.console-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center;
  transition: transform 0.3s;
}

.console-card:hover .console-img {
  transform: scale(1.05);
}

.quick-view {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(0, 102, 204, 0.9);
  color: white;
  text-align: center;
  padding: 10px;
  transform: translateY(100%);
  transition: transform 0.3s;
  font-size: 0.9rem;
  font-weight: 500;
}

.console-image:hover .quick-view {
  transform: translateY(0);
}

.console-type {
  position: absolute;
  top: 1rem;
  right: 1rem;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: bold;
  backdrop-filter: blur(10px);
  z-index: 1;
}

.console-type.ps5 {
  background: rgba(0, 102, 204, 0.9);
  color: white;
}

.console-type.ps4 {
  background: rgba(0, 87, 184, 0.9);
  color: white;
}

.console-type.xbox {
  background: rgba(16, 124, 16, 0.9);
  color: white;
}

.console-info {
  padding: 1.5rem;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.clickable-title {
  cursor: pointer;
  transition: color 0.3s;
}

.clickable-title:hover {
  color: #0066cc;
}

.console-info h3 {
  font-size: 1.3rem;
  margin-bottom: 0.5rem;
  color: #2d3748;
  min-height: 3rem;
}

.price {
  font-size: 1.8rem;
  font-weight: bold;
  color: #0066cc;
  margin: 0.5rem 0;
}

.price span {
  font-size: 0.9rem;
  color: #718096;
  font-weight: normal;
}

.description {
  color: #718096;
  margin: 0.5rem 0;
  font-size: 0.9rem;
  line-height: 1.4;
  flex-grow: 1;
}

.availability {
  margin: 1rem 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 500;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.availability .status-dot {
  background: #e53e3e;
}

.availability.available .status-dot {
  background: #38a169;
}

.availability.available {
  color: #38a169;
}

.button-group {
  display: flex;
  gap: 1rem;
  margin-top: auto;
}

.details-btn {
  flex: 1;
  padding: 0.85rem;
  background: white;
  color: #0066cc;
  border: 2px solid #0066cc;
  border-radius: 12px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
  transition: all 0.3s;
}

.details-btn:hover {
  background: #0066cc;
  color: white;
}

.rent-btn {
  flex: 1;
  padding: 0.85rem;
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  color: white;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
  transition: all 0.3s;
}

.rent-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 102, 204, 0.3);
}

/* Модальное окно деталей */
/* Модальное окно деталей */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(5px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-content {
  background: white;
  border-radius: 20px;
  position: relative;
  max-height: 90vh;
  overflow-y: auto;
  animation: slideIn 0.3s ease;
}

.details-modal {
  max-width: 1000px;
  width: 100%;
}

.modal-close {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  border: none;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  cursor: pointer;
  font-size: 1.2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
  transition: all 0.3s;
}

.modal-close:hover {
  background: rgba(0, 0, 0, 0.8);
  transform: scale(1.1);
}

.details-container {
  display: grid;
  grid-template-columns: 1fr 1.5fr;
  gap: 2rem;
  padding: 2rem;
}

.details-image {
  position: relative;
  border-radius: 15px;
  overflow: hidden;
}

.details-image img {
  width: 100%;
  height: auto;
  object-fit: cover;
}

.console-badge {
  position: absolute;
  top: 1rem;
  right: 1rem;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-weight: bold;
}

.details-info h2 {
  font-size: 1.8rem;
  color: #2d3748;
  margin-bottom: 1rem;
}

.details-price {
  background: linear-gradient(135deg, #f0f7ff 0%, #e8f0fe 100%);
  padding: 1rem;
  border-radius: 10px;
  margin-bottom: 1.5rem;
}

.details-price .price-label {
  color: #4a5568;
  margin-right: 1rem;
}

.details-price .price-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: #0066cc;
}

.details-section {
  margin-bottom: 1.5rem;
}

.details-section h3 {
  color: #2d3748;
  margin-bottom: 0.75rem;
  font-size: 1.2rem;
}

.details-section ul {
  list-style: none;
  padding-left: 0;
}

.details-section li {
  padding: 0.5rem 0;
  color: #4a5568;
}

.games-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.game-tag {
  background: #f7fafc;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.9rem;
  color: #2d3748;
  border: 1px solid #e2e8f0;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}

.feature {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  background: #f7fafc;
  border-radius: 10px;
}

.feature-icon {
  font-size: 1.5rem;
}

.reviews {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.review {
  background: #f7fafc;
  padding: 1rem;
  border-radius: 10px;
}

.review-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.reviewer {
  font-weight: bold;
  color: #2d3748;
}

.rating {
  color: #fbbf24;
}

/* Стили для кнопок в модальном окне деталей */
.details-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}

.details-actions .rent-now-btn {
  flex: 2;
  padding: 0.85rem;
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  color: white;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
  transition: all 0.3s;
}

.details-actions .rent-now-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 102, 204, 0.3);
}

.details-actions .btn-secondary {
  flex: 1;
  padding: 0.85rem;
  background: #f7fafc;
  color: #4a5568;
  border: 2px solid #e2e8f0;
  border-radius: 12px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
  transition: all 0.3s;
}

.details-actions .btn-secondary:hover {
  background: #edf2f7;
  transform: translateY(-2px);
}

@keyframes slideIn {
  from {
    transform: translateY(-50px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

/* Адаптивность */
@media (max-width: 768px) {
  .home {
    padding: 1rem;
  }

  .hero {
    padding: 2rem 1rem;
  }

  .hero h1 {
    font-size: 2rem;
  }

  .details-container {
    grid-template-columns: 1fr;
    gap: 1rem;
    padding: 1rem;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .button-group {
    flex-direction: column;
  }

  .details-actions {
    flex-direction: column;
    gap: 0.75rem;
  }

  .details-actions .rent-now-btn,
  .details-actions .btn-secondary {
    width: 100%;
  }

  .modal-content {
    width: 95%;
    padding: 1rem;
  }

  .details-modal {
    max-width: 95%;
  }
}

@media (max-width: 480px) {
  .hero h1 {
    font-size: 1.5rem;
  }

  .price {
    font-size: 1.3rem;
  }

  .console-info h3 {
    font-size: 1rem;
  }

  .description {
    font-size: 0.85rem;
  }

  .details-info h2 {
    font-size: 1.3rem;
  }

  .details-price .price-value {
    font-size: 1.2rem;
  }

  .game-tag {
    font-size: 0.8rem;
    padding: 0.35rem 0.75rem;
  }

  .feature {
    font-size: 0.9rem;
  }

  .review p {
    font-size: 0.85rem;
  }
}
</style>
