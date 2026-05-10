<template>
  <div class="rentals-page">
    <div class="page-header">
      <h1>Мои аренды</h1>
      <p>История и статус ваших аренд</p>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="loader"></div>
      <p>Загрузка ваших аренд...</p>
    </div>

    <div v-else-if="!rentals || rentals.length === 0" class="no-rentals">
      <div class="empty-state">
        <div class="empty-icon">📦</div>
        <h3>У вас пока нет активных аренд</h3>
        <p>Выберите консоль и начните играть уже сегодня!</p>
        <router-link to="/" class="browse-btn"
          >Посмотреть консоли →</router-link
        >
      </div>
    </div>

    <div v-else class="rentals-list">
      <div v-for="rental in rentals" :key="rental.id" class="rental-card">
        <div class="rental-image">
          <img
            :src="rental.console?.image_url || '/placeholder.png'"
            :alt="rental.console?.model || 'Консоль'"
            @error="handleImageError"
            class="console-img"
          />
          <div
            class="console-badge"
            :class="getConsoleClass(rental.console?.type)"
          >
            {{ rental.console?.type || "Консоль" }}
          </div>
        </div>
        <div class="rental-details">
          <h3>{{ rental.console?.model || "Консоль" }}</h3>
          <div class="rental-info">
            <div class="info-item">
              <span class="info-label">📅 Начало:</span>
              <span class="info-value">{{
                formatDate(rental.start_date)
              }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">📅 Окончание:</span>
              <span class="info-value">{{ formatDate(rental.end_date) }}</span>
            </div>
          </div>
          <div class="rental-price">
            <span class="price-label">💰 Стоимость:</span>
            <span class="price-value">{{ rental.total_price }} ₽</span>
          </div>
          <div v-if="rental.delivery_address" class="delivery-info">
            <span class="info-label">🚚 Адрес доставки:</span>
            <span class="info-value">{{ rental.delivery_address }}</span>
          </div>
          <div class="status-badge" :class="rental.status">
            <span class="status-icon">{{ getStatusIcon(rental.status) }}</span>
            {{ getStatusText(rental.status) }}
          </div>
        </div>
        <button
          v-if="rental.status === 'active'"
          @click="returnConsole(rental.id)"
          class="return-btn"
        >
          Вернуть консоль
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";

export default {
  name: "MyRentals",
  setup() {
    const router = useRouter();
    const rentals = ref([]);
    const loading = ref(true);
    let isComponentMounted = true;

    const fetchRentals = async () => {
      // Проверяем наличие токена
      const token = localStorage.getItem("auth_token");
      if (!token) {
        router.push("/");
        return;
      }

      loading.value = true;
      try {
        const response = await axios.get("/api/my-rentals");
        if (isComponentMounted) {
          rentals.value = Array.isArray(response.data) ? response.data : [];
        }
      } catch (error) {
        console.error("Error fetching rentals:", error);
        if (error.response?.status === 401) {
          // Если не авторизован, перенаправляем на главную
          localStorage.removeItem("auth_token");
          router.push("/");
        }
        if (isComponentMounted) {
          rentals.value = [];
        }
      } finally {
        if (isComponentMounted) {
          loading.value = false;
        }
      }
    };

    const returnConsole = async (id) => {
      if (confirm("Вы уверены, что хотите вернуть консоль?")) {
        try {
          await axios.put(`/api/rentals/${id}/return`);
          alert("Консоль успешно возвращена");
          await fetchRentals();
        } catch (error) {
          console.error("Error returning console:", error);
          if (error.response?.status === 401) {
            localStorage.removeItem("auth_token");
            router.push("/");
          } else {
            alert("Ошибка при возврате");
          }
        }
      }
    };

    const formatDate = (date) => {
      if (!date) return "Дата не указана";
      try {
        return new Date(date).toLocaleDateString("ru-RU", {
          day: "numeric",
          month: "long",
          year: "numeric",
        });
      } catch (e) {
        return "Дата не указана";
      }
    };

    const getStatusText = (status) => {
      const statusMap = {
        active: "Активна",
        returned: "Возвращена",
        cancelled: "Отменена",
      };
      return statusMap[status] || status || "Неизвестно";
    };

    const getStatusIcon = (status) => {
      const iconMap = {
        active: "🎮",
        returned: "✅",
        cancelled: "❌",
      };
      return iconMap[status] || "📦";
    };

    const getConsoleClass = (type) => {
      if (type === "PS5") return "ps5";
      if (type === "PS4") return "ps4";
      if (type === "XBOX") return "xbox";
      return "";
    };

    const handleImageError = (event) => {
      event.target.src = "https://via.placeholder.com/150x150?text=Console";
    };

    onMounted(() => {
      isComponentMounted = true;
      fetchRentals();
    });

    onBeforeUnmount(() => {
      isComponentMounted = false;
    });

    return {
      rentals,
      loading,
      returnConsole,
      formatDate,
      getStatusText,
      getStatusIcon,
      getConsoleClass,
      handleImageError,
    };
  },
};
</script>
<style scoped>
.rentals-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
}

.page-header {
  text-align: center;
  margin-bottom: 3rem;
}

.page-header h1 {
  font-size: 2.5rem;
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  margin-bottom: 0.5rem;
}

.page-header p {
  color: #718096;
  font-size: 1.1rem;
}

.loading-state {
  text-align: center;
  padding: 4rem;
}

.loader {
  border: 3px solid #f3f3f3;
  border-top: 3px solid #0066cc;
  border-radius: 50%;
  width: 50px;
  height: 50px;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.no-rentals {
  background: white;
  border-radius: 20px;
  padding: 3rem;
  text-align: center;
  border: 1px solid rgba(0, 102, 204, 0.1);
}

.empty-state {
  max-width: 400px;
  margin: 0 auto;
}

.empty-icon {
  font-size: 5rem;
  margin-bottom: 1rem;
}

.empty-state h3 {
  color: #2d3748;
  margin-bottom: 0.5rem;
}

.empty-state p {
  color: #718096;
  margin-bottom: 2rem;
}

.browse-btn {
  display: inline-block;
  padding: 0.85rem 2rem;
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  color: white;
  text-decoration: none;
  border-radius: 12px;
  font-weight: 600;
  transition: all 0.3s;
}

.browse-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 102, 204, 0.3);
}

.rentals-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.rental-card {
  background: white;
  border-radius: 20px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1.5rem;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
  transition: all 0.3s;
  border: 1px solid rgba(0, 102, 204, 0.1);
}

.rental-card:hover {
  transform: translateX(5px);
  box-shadow: 0 10px 30px rgba(0, 102, 204, 0.1);
}

.rental-image {
  position: relative;
  width: 120px;
  height: 120px;
  flex-shrink: 0;
  border-radius: 15px;
  overflow: hidden;
  background: linear-gradient(135deg, #f7fafc 0%, #edf2f7 100%);
}

.console-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.rental-card:hover .console-img {
  transform: scale(1.05);
}

.console-badge {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  padding: 0.25rem 0.75rem;
  border-radius: 15px;
  font-size: 0.75rem;
  font-weight: bold;
  backdrop-filter: blur(5px);
}

.console-badge.ps5 {
  background: rgba(0, 102, 204, 0.9);
  color: white;
}

.console-badge.ps4 {
  background: rgba(0, 87, 184, 0.9);
  color: white;
}

.console-badge.xbox {
  background: rgba(16, 124, 16, 0.9);
  color: white;
}

.rental-details {
  flex: 1;
}

.rental-details h3 {
  color: #2d3748;
  margin-bottom: 0.75rem;
  font-size: 1.2rem;
}

.rental-info {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 0.5rem;
  flex-wrap: wrap;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #718096;
  font-size: 0.9rem;
}

.info-label {
  font-weight: 500;
  color: #4a5568;
}

.info-value {
  color: #2d3748;
}

.rental-price {
  margin: 0.5rem 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.price-label {
  color: #718096;
}

.price-value {
  font-weight: bold;
  color: #0066cc;
  font-size: 1.1rem;
}

.delivery-info {
  margin: 0.5rem 0;
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.35rem 1rem;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 500;
  margin-top: 0.5rem;
  width: fit-content;
}

.status-icon {
  font-size: 1rem;
}

.status-badge.active {
  background: #c6f6d5;
  color: #276749;
}

.status-badge.returned {
  background: #e2e8f0;
  color: #4a5568;
}

.status-badge.cancelled {
  background: #fed7d7;
  color: #c53030;
}

.return-btn {
  padding: 0.7rem 1.5rem;
  background: white;
  color: #e53e3e;
  border: 2px solid #e53e3e;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s;
  flex-shrink: 0;
}

.return-btn:hover {
  background: #e53e3e;
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(229, 62, 62, 0.3);
}

@media (max-width: 768px) {
  .rentals-page {
    padding: 1rem;
  }

  .page-header h1 {
    font-size: 1.8rem;
  }

  .rental-card {
    flex-direction: column;
    text-align: center;
    padding: 1rem;
  }

  .rental-image {
    width: 150px;
    height: 150px;
    margin: 0 auto;
  }

  .rental-info {
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
  }

  .rental-details h3 {
    text-align: center;
  }

  .delivery-info {
    justify-content: center;
  }

  .return-btn {
    width: 100%;
  }

  .status-badge {
    margin: 0.5rem auto;
  }
}

@media (max-width: 480px) {
  .rental-image {
    width: 120px;
    height: 120px;
  }

  .info-item {
    flex-direction: column;
    gap: 0.25rem;
  }

  .price-value {
    font-size: 1rem;
  }
}
</style>
