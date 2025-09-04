import api from '@/config/api';
import { OrganizationRequest, OrganizationResponse } from '@/models/Organization';

export const OrganizationService = {
  /**
   * Crée une nouvelle organisation
   * @param organizationData Les données de l'organisation à créer
   * @returns La réponse du serveur avec l'ID de l'organisation créée
   */
  async createOrganization(organizationData: OrganizationRequest): Promise<OrganizationResponse> {
    try {
      const response = await fetch(api.ORGANIZATION, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(organizationData)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Erreur lors de la création de l\'organisation');
      }

      return await response.json();
    } catch (error) {
      console.error('Error creating organization:', error);
      throw error;
    }
  },

  /**
   * Récupère les détails d'une organisation par son ID
   * @param organizationId L'ID de l'organisation à récupérer
   * @returns Les détails de l'organisation
   */
  async getOrganization(organizationId: string): Promise<OrganizationRequest> {
    try {
      const response = await fetch(`${api.ORGANIZATION}/${organizationId}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      if (!response.ok) {
        throw new Error('Erreur lors de la récupération de l\'organisation');
      }

      return await response.json();
    } catch (error) {
      console.error('Error fetching organization:', error);
      throw error;
    }
  },

  /**
   * Met à jour une organisation existante
   * @param organizationId L'ID de l'organisation à mettre à jour
   * @param organizationData Les nouvelles données de l'organisation
   * @returns La réponse du serveur
   */
  async updateOrganization(organizationId: string, organizationData: Partial<OrganizationRequest>): Promise<OrganizationResponse> {
    try {
      const response = await fetch(`${api.ORGANIZATION}/${organizationId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(organizationData)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Erreur lors de la mise à jour de l\'organisation');
      }

      return await response.json();
    } catch (error) {
      console.error('Error updating organization:', error);
      throw error;
    }
  }
};
